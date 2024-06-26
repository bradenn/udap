// Copyright (c) 2022 Braden Nicholson

package services

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
	"udap/internal/log"
	"udap/internal/pulse"
)

func NewModuleService(repository ports.ModuleRepository, runtime ports.ModuleOperator) ports.ModuleService {
	return &moduleService{
		repository:    repository,
		operator:      runtime,
		timeoutsMutex: sync.RWMutex{},
		timeouts:      map[string]*time.Timer{},
	}
}

const DIR = "modules"

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

type moduleService struct {
	repository    ports.ModuleRepository
	operator      ports.ModuleOperator
	timeouts      map[string]*time.Timer
	timeoutsMutex sync.RWMutex
	generic.Watchable[domain.Module]
}

func (u *moduleService) HandleEmits(mutation domain.Mutation) error {
	err := u.operator.HandleEmit(mutation)
	if err != nil {
		return err
	}
	return nil
}

func (u *moduleService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	for _, module := range *all {
		err = u.Emit(module)
		if err != nil {
			return err
		}
	}
	return nil
}

// Panic handles a panic state
func (u *moduleService) Panic(module domain.Module) error {
	log.Recovered("Module '%s' panicked; module entering safe-mode", module.Name)
	err := u.Dispose(module.Id)
	if err != nil {
		log.Err(fmt.Errorf("module disposal failed, runtime must be flushed to resume operation: %s",
			err.Error()))
		return err
	}
	return nil
}

func (u *moduleService) Update(id string) error {
	// Find the local module from the repository
	module, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	// Make sure the module is valid
	if !module.Enabled || !module.Running {
		return fmt.Errorf("module must be enabled and running to update")
	}
	// Catch any panics that may occur when running the update function
	defer func() {
		// Attempt to recover from a panic
		if r := recover(); r != nil {
			err = u.Panic(*module)
			if err != nil {
				return
			}
		}
	}()
	addr := fmt.Sprintf("module.%s.update", module.UUID)
	// Mark the time that the update begins
	pulse.Begin(addr)
	// End the pulse when the update concludes or errors out
	defer pulse.End(addr)
	// Attempt to update the modules
	err = u.operator.Update(module.UUID)
	if err != nil {
		return err
	}

	go func(mod *domain.Module) {
		duration := mod.Interval
		if duration < time.Second {
			return
		}
		timer := time.NewTimer(duration)
		u.timeoutsMutex.Lock()
		u.timeouts[mod.UUID] = timer
		u.timeoutsMutex.Unlock()
		select {
		case _ = <-timer.C:
			err = u.Update(id)
			if err != nil {
				return
			}
		}
	}(module)
	// Return normally
	return nil
}

// Run runs the startup code for each module, not to be confused with the setup function with connects and the module
func (u *moduleService) Run(id string) error {
	module, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	if module.Running {
		return fmt.Errorf("module must not be running to run")
	}
	// Mark the module start as starting
	err = u.setState(module.Id, STARTING)
	if err != nil {
		return err
	}
	// Catch any panics that may occur when running the run function
	defer func() {
		// Attempt to recover from a panic
		if r := recover(); r != nil {
			err = u.Panic(*module)
			if err != nil {
				return
			}
		}
	}()
	addr := fmt.Sprintf("module.%s.run", module.UUID)
	// Mark the time that the update begins
	pulse.Begin(addr)
	// End the pulse when the update concludes or errors out
	defer pulse.End(addr)
	// Attempt to run the module
	err = u.operator.Run(module.UUID)
	if err != nil {
		// Set the module state to error if the run fails
		err = u.setState(module.Id, ERROR)
		if err != nil {
			return err
		}
		return err
	}
	// Set the module as running so it can begin updating
	module.Running = true
	err = u.repository.Update(module)
	if err != nil {
		return err
	}
	// Mark the module as running
	err = u.setState(module.Id, RUNNING)
	if err != nil {
		return err
	}
	if module.Enabled && module.Running {
		err = u.Update(module.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *moduleService) Load(id string) error {
	module, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	// Catch any panics that may occur when running the run function
	defer func() {
		// Attempt to recover from a panic
		if r := recover(); r != nil {
			err = u.Panic(*module)
			if err != nil {
				return
			}
		}
	}()
	start := time.Now()
	// Attempt to load the module
	config, err := u.operator.Load(module.Name, module.UUID)
	if err != nil {
		return err
	}
	module.Variables = config.Variables
	module.Version = config.Version
	module.Description = config.Description
	module.Interval = config.Interval
	module.Type = config.Type
	module.Author = config.Author
	module.Running = false
	log.Event("Module '%s' @ 0x%s loaded. (%s)", module.Name, module.SessionId(),
		time.Since(start).String())
	err = u.repository.Update(module)
	if err != nil {
		return err
	}
	// Set the state if the operation was a success, setState will also update the module data
	err = u.setState(module.Id, IDLE)
	if err != nil {
		return err
	}

	return nil
}

func (u *moduleService) Build(id string) error {
	module, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	uid := uuid.New().String()
	module.UUID = uid
	err = u.repository.Update(module)
	if err != nil {
		return err
	}
	// start := time.Now()
	err = u.operator.Build(module.Name, module.UUID)
	if err != nil {
		return err
	}
	//log.Event("Module '%s' @ 0x%s compiled. (%s)", module.Name, module.SessionId(),
	//	time.Since(start).Truncate(time.Millisecond).String())
	return nil
}

// Dispose halts a modules activity and destroys its runtime
func (u *moduleService) Dispose(id string) error {
	module, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	if !module.Enabled || !module.Running {
		return fmt.Errorf("module must be enabled and running to dispose")
	}
	// Mark the state as halting, since it may take a while
	err = u.setState(module.Id, HALTING)
	if err != nil {
		return err
	}
	start := time.Now()
	// Attempt to dispose of the module (only works if the module developer plays nicely)
	err = u.operator.Dispose(module.Name, module.UUID)
	if err != nil {
		return err
	}
	// Set the module as not running, so it is not updated
	module.Running = false
	log.Event("Module '%s' @ 0x%s unloaded. (%s)", module.Name, module.SessionId(),
		time.Since(start).Truncate(time.Millisecond).String())
	module.UUID = ""
	// Mark the module as stopped if the disposal was successful
	err = u.setState(module.Id, STOPPED)
	if err != nil {
		return err
	}

	return nil
}

func (u *moduleService) UpdateAll() error {
	//modules, err := u.repository.FindAll()
	//if err != nil {
	//	return err
	//}
	//wg := sync.WaitGroup{}
	//ref := *modules
	//for _, module := range ref {
	//	if !module.Running || !module.Enabled {
	//		continue
	//	}
	//	wg.Add(1)
	//
	//	go func(mod domain.Module) {
	//		defer wg.Done()
	//		err = u.Update(mod.Id)
	//		if err != nil {
	//			log.Err(err)
	//		}
	//	}(module)
	//}
	//wg.Wait()
	return nil
}

func (u *moduleService) setState(id string, state string) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	byId.State = state
	err = u.repository.Update(byId)
	if err != nil {
		return err
	}
	err = u.Emit(*byId)
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
	wg := sync.WaitGroup{}
	wg.Add(len(*modules))
	for _, module := range *modules {
		go func(mod domain.Module) {
			defer wg.Done()
			if !mod.Enabled {
				return
			}
			err = u.Run(mod.Id)
			if err != nil {
				log.Err(err)
				return
			}
		}(module)
	}
	wg.Wait()
	return nil
}

func (u *moduleService) DisposeAll() error {
	modules, err := u.repository.FindAll()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(*modules))
	for _, module := range *modules {
		go func(mod domain.Module) {
			defer wg.Done()
			if !mod.Enabled || !mod.Running {
				return
			}
			err = u.Dispose(mod.Id)
			if err != nil {
				log.Err(err)
				return
			}
		}(module)
	}
	wg.Wait()
	return nil
}

func (u *moduleService) Loader(wg *sync.WaitGroup) chan domain.Module {
	threads := 8
	channel := make(chan domain.Module, threads)

	for i := 0; i < threads; i++ {
		go func() {
			for module := range channel {

				err := u.Build(module.Id)
				if err != nil {
					wg.Done()
					log.Err(err)
					continue
				}
				log.Event("Loaded module '%s'", module.Name)
				wg.Done()
			}
		}()
	}

	return channel
}

func (u *moduleService) LoadAll() error {
	modules, err := u.repository.FindAll()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(*modules))
	//channel := u.Loader(&wg)
	for _, module := range *modules {
		//channel <- module
		go func(mod domain.Module) {
			defer wg.Done()
			err = u.Load(mod.Id)
			if err != nil {
				log.Err(err)
				return
			}
		}(module)
	}

	wg.Wait()

	log.Event("All modules loaded")

	return nil
}

func (u *moduleService) Discover() error {
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
				continue
			}
		}
	}
	return nil
}

func (u *moduleService) Builder(wg *sync.WaitGroup) chan domain.Module {
	threads := runtime.NumCPU()
	channel := make(chan domain.Module, threads)

	for i := 0; i < threads; i++ {
		go func() {
			for module := range channel {
				start := time.Now()

				err := u.Build(module.Id)
				if err != nil {
					wg.Done()
					log.Err(err)
					err = u.setState(module.Id, ERROR)
					continue
				}

				err = u.setState(module.Id, UNINITIALIZED)
				if err != nil {
					continue
				}

				log.Event("Module '%s' @ 0x%s compiled. (%s)", module.Name, module.SessionId(),
					time.Since(start).String())

				wg.Done()
			}
		}()
	}

	return channel
}

func (u *moduleService) BuildAll() error {
	modules, err := u.repository.FindAll()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(*modules))
	builder := u.Builder(&wg)
	for _, module := range *modules {
		module.Running = false
		builder <- module
	}
	wg.Wait()
	close(builder)
	return nil
}

func (u *moduleService) GetConfig(id string, key string) (string, error) {
	byId, err := u.repository.FindByUUID(id)
	if err != nil {
		return "", err
	}

	var values map[string]string
	err = json.Unmarshal([]byte(byId.Config), &values)
	if err != nil {
		err = nil
		values = map[string]string{}
	}

	value, ok := values[key]
	if !ok {
		return "", fmt.Errorf("config value on key '%s' does not exist", key)
	}

	return value, nil
}

func (u *moduleService) InitConfig(id string, key string, value string) error {
	byId, err := u.repository.FindByUUID(id)
	if err != nil {
		return err
	}

	var values map[string]string
	err = json.Unmarshal([]byte(byId.Config), &values)
	if err != nil {
		err = nil
		values = map[string]string{}
	}

	val, ok := values[key]
	if val != "" && ok {
		return nil
	}

	values[key] = value

	marshal, err := json.Marshal(values)
	if err != nil {
		return err
	}

	byId.Config = string(marshal)

	err = u.save(byId)
	if err != nil {
		return err
	}
	return nil
}

func (u *moduleService) SetConfig(id string, key string, value string) error {
	byId, err := u.repository.FindByUUID(id)
	if err != nil {
		return err
	}

	var values map[string]string
	err = json.Unmarshal([]byte(byId.Config), &values)
	if err != nil {
		err = nil
		values = map[string]string{}
	}

	values[key] = value

	marshal, err := json.Marshal(values)
	if err != nil {
		return err
	}

	byId.Config = string(marshal)

	err = u.save(byId)
	if err != nil {
		return err
	}

	return nil
}

// Repository Mapping

func (u *moduleService) FindAll() (*[]domain.Module, error) {
	return u.repository.FindAll()
}

func (u *moduleService) FindByName(name string) (*domain.Module, error) {
	return u.repository.FindByName(name)
}

func (u *moduleService) Disable(id string) error {
	module, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	module.Enabled = false
	err = u.repository.Update(module)
	if err != nil {
		return err
	}
	err = u.Emit(*module)
	if err != nil {
		return err
	}
	return nil
}

func (u *moduleService) save(module *domain.Module) error {
	err := u.repository.Update(module)
	if err != nil {
		return err
	}
	err = u.Emit(*module)
	if err != nil {
		return err
	}
	return nil
}

func (u *moduleService) Enable(id string) error {
	module, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	module.Enabled = true
	err = u.repository.Update(module)
	if err != nil {
		return err
	}
	err = u.Emit(*module)
	if err != nil {
		return err
	}
	return nil
}

func (u *moduleService) Reload(name string) error {

	target, err := u.FindByName(name)
	if err != nil {
		return err
	}
	if !target.Enabled {
		return fmt.Errorf("module must be enabled and running to reload")
	}
	if target.Running {
		err = u.Dispose(target.Id)
		if err != nil {
			return err
		}
	}
	err = u.Build(target.Id)
	if err != nil {
		return err
	}
	err = u.Load(target.Id)
	if err != nil {
		return err
	}
	err = u.Run(target.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *moduleService) Halt(name string) error {
	byName, err := u.FindByName(name)
	if err != nil {
		return err
	}
	err = u.Dispose(byName.Id)
	if err != nil {
		return err
	}
	return nil
}
