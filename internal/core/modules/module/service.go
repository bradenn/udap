// Copyright (c) 2022 Braden Nicholson

package module

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"udap/internal/core/domain"
	"udap/internal/log"
)

const DIR = "modules"

type moduleService struct {
	repository domain.ModuleRepository
	operator   domain.ModuleOperator
	channel    chan<- domain.Mutation
}

const (
	DISCOVERED    = "discovered"
	UNINITIALIZED = "uninitialized"
	IDLE          = "idle"
	STARTING      = "starting"
	RUNNING       = "running"
	HALTING       = "halting"
	STOPPED       = "stopped"
	ERROR         = "error"
)

func (u *moduleService) Watch(ref chan<- domain.Mutation) error {
	if u.channel != nil {
		return fmt.Errorf("channel in use")
	}
	u.channel = ref
	return nil
}

func (u *moduleService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	for _, module := range *all {
		err = u.emit(&module)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *moduleService) emit(module *domain.Module) error {
	if u.channel == nil {
		return fmt.Errorf("channel is null")
	}
	u.channel <- domain.Mutation{
		Status:    "update",
		Operation: "module",
		Body:      *module,
		Id:        module.Id,
	}
	return nil
}

func (u *moduleService) Update(module *domain.Module) error {
	return u.operator.Update(module)
}

// Run runs the startup code for each module, not to be confused with the setup function with connects and the module
func (u *moduleService) Run(module *domain.Module) error {
	// Mark the module start as starting
	err := u.setState(module, STARTING)
	if err != nil {
		return err
	}
	// Attempt to run the module
	err = u.operator.Run(module)
	if err != nil {
		// Set the module state to error if the run fails
		err = u.setState(module, ERROR)
		if err != nil {
			return err
		}
		return err
	}
	// Set the module as running so it can begin updating
	module.Running = false
	// Mark the module as running
	err = u.setState(module, RUNNING)
	if err != nil {
		return err
	}
	return nil
}

func (u *moduleService) Load(module *domain.Module) error {
	// Attempt to load the module
	err := u.operator.Load(module)
	if err != nil {
		return err
	}
	// Set the state if the operation was a success, setState will also update the module data
	err = u.setState(module, IDLE)
	if err != nil {
		return err
	}
	return nil
}

func (u *moduleService) Build(module *domain.Module) error {
	return u.operator.Build(module)
}

// Dispose halts a modules activity and destroys its runtime
func (u *moduleService) Dispose(module *domain.Module) error {
	// Mark the state as halting, since it may take a while
	err := u.setState(module, HALTING)
	if err != nil {
		return err
	}
	// Attempt to dispose of the module (only works if the module developer plays nicely)
	err = u.operator.Dispose(module)
	if err != nil {
		return err
	}
	// Set the module as not running, so it is not updated
	module.Running = false
	// Mark the module as stopped if the disposal was successful
	err = u.setState(module, STOPPED)
	if err != nil {
		return err
	}
	return nil
}

func (u *moduleService) UpdateAll() error {
	modules, err := u.repository.FindAll()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	ref := *modules
	wg.Add(len(ref))
	for _, module := range ref {
		go func(mod domain.Module) {
			defer wg.Done()
			err = u.Update(&mod)
			if err != nil {
				log.Err(err)
			}
		}(module)
	}
	wg.Wait()
	return nil
}

func (u *moduleService) setState(module *domain.Module, state string) error {
	module.State = state
	err := u.repository.Update(module)
	if err != nil {
		return err
	}
	err = u.emit(module)
	if err != nil {
		return err
	}
	return nil
}

func (u *moduleService) RunAll() error {
	modules, err := u.repository.FindAll()
	if err != nil {
		return err
	}

	for _, module := range *modules {
		go func(mod domain.Module) {
			err = u.Run(&mod)
			if err != nil {
				log.Err(err)
				return
			}
		}(module)
	}

	return nil
}

func (u *moduleService) LoadAll() error {
	modules, err := u.repository.FindAll()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(*modules))
	for _, module := range *modules {
		go func(mod domain.Module) {
			defer wg.Done()
			err = u.Load(&mod)
			if err != nil {
				log.Err(err)
				return
			}
		}(module)
	}
	wg.Wait()
	return nil
}

func NewService(repository domain.ModuleRepository, operator domain.ModuleOperator) domain.ModuleService {
	return &moduleService{
		repository: repository,
		operator:   operator,
	}
}

func (u moduleService) Discover() error {
	// Format the pattern for glob search
	pattern := fmt.Sprintf("./%s/*/*.go", DIR)
	// Run the search for go files
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	// Launch a go func to build each one
	for _, p := range files {
		name := strings.Replace(filepath.Base(p), ".go", "", 1)
		var target *domain.Module
		target, err = u.repository.FindByName(name)
		if err != nil {
			target = &domain.Module{}
			target.Name = name
			target.Path = p
			target.State = DISCOVERED
			err = u.repository.Create(target)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *moduleService) BuildAll() error {
	modules, err := u.repository.FindAll()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(*modules))
	for _, module := range *modules {
		go func(mod domain.Module) {
			defer wg.Done()
			err = u.Build(&mod)
			if err != nil {
				log.Err(err)
				err = u.setState(&mod, ERROR)
				if err != nil {
					return
				}
				return
			}
			err = u.setState(&mod, UNINITIALIZED)
			if err != nil {
				return
			}
		}(module)
	}
	wg.Wait()
	return nil
}

// Repository Mapping

func (u moduleService) FindAll() (*[]domain.Module, error) {
	return u.repository.FindAll()
}

func (u moduleService) FindByName(name string) (*domain.Module, error) {
	return u.repository.FindByName(name)
}

func (u moduleService) Disable(id string) error {
	module, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	module.Enabled = false
	err = u.repository.Update(module)
	if err != nil {
		return err
	}
	err = u.emit(module)
	if err != nil {
		return err
	}
	return nil
}

func (u moduleService) save(module *domain.Module) error {
	err := u.repository.Update(module)
	if err != nil {
		return err
	}
	err = u.emit(module)
	if err != nil {
		return err
	}
	return nil
}

func (u moduleService) Enable(id string) error {
	module, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	module.Enabled = true
	err = u.repository.Update(module)
	if err != nil {
		return err
	}
	err = u.emit(module)
	if err != nil {
		return err
	}
	return nil
}

func (u moduleService) Reload(name string) error {
	return nil
}

func (u moduleService) Halt(name string) error {
	byName, err := u.FindByName(name)
	if err != nil {
		return err
	}
	err = u.Dispose(byName)
	if err != nil {
		return err
	}
	return nil
}
