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
}

func (u *moduleService) Update(module *domain.Module) error {
	return u.operator.Update(module)
}

func (u *moduleService) Run(module *domain.Module) error {
	return u.operator.Run(module)
}

func (u *moduleService) Load(module *domain.Module) error {
	err := u.operator.Load(module)
	if err != nil {
		return err
	}
	err = u.repository.Update(module)
	if err != nil {
		return err
	}
	return nil
}

func (u *moduleService) Build(module *domain.Module) error {
	return u.operator.Build(module)
}

func (u *moduleService) UpdateAll() error {
	modules, err := u.repository.FindAll()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(*modules))
	for _, module := range *modules {
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
			target.Name = name
			target.Path = p
			err = u.repository.Create(target)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (u moduleService) BuildAll() error {
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

func (u moduleService) Disable(name string) error {
	_, err := u.FindByName(name)
	if err != nil {
		return err
	}
	return nil
}

func (u moduleService) Enable(name string) error {
	// TODO implement me
	panic("implement me")
}

func (u moduleService) Reload(name string) error {
	// TODO implement me
	panic("implement me")
}

func (u moduleService) Halt(name string) error {
	// TODO implement me
	panic("implement me")
}
