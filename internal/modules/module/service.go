// Copyright (c) 2022 Braden Nicholson

package module

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
)

const DIR = "modules"

type moduleService struct {
	repository domain.ModuleRepository
}

func NewService(repository domain.ModuleRepository) domain.ModuleService {
	return moduleService{repository: repository}
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

func (u moduleService) Build(module *domain.Module) error {
	start := time.Now()
	if _, err := os.Stat(module.Path); err != nil {
		return err
	}
	// Create a timeout to prevent modules from taking too long to build
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
	// Cancel the timeout of it exits before the timeout is up
	defer cancelFunc()
	binary := strings.Replace(module.Path, ".go", ".so", 1)
	// Prepare the command arguments
	args := []string{"build", "-v", "-buildmode=plugin", "-o", binary, module.Path}
	// Initialize the command structure
	cmd := exec.CommandContext(timeout, "go", args...)
	// Run and get the stdout and stderr from the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.ErrF(errors.New(string(output)), "Module '%s' build failed:", module.Name)
		return nil
	}
	log.Event("Module '%s' compiled successfully (%s)", module.Name, time.Since(start).Truncate(time.Millisecond).String())
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
