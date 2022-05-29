// Copyright (c) 2022 Braden Nicholson

package module

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/pkg/plugin"
)

type moduleOperator struct {
	ctrl    *controller.Controller
	modules map[string]plugin.ModuleInterface
}

func (m moduleOperator) Update(module *domain.Module) error {
	err := m.modules[module.Name].Update()
	if err != nil {
		return err
	}
	return nil
}

func (m moduleOperator) Run(module *domain.Module) error {
	err := m.modules[module.Name].Run()
	if err != nil {
		return err
	}
	return nil
}

func NewOperator(ctrl *controller.Controller) domain.ModuleOperator {
	return &moduleOperator{
		ctrl:    ctrl,
		modules: map[string]plugin.ModuleInterface{},
	}
}

func (m moduleOperator) Build(module *domain.Module) error {
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

func (m moduleOperator) Load(module *domain.Module) error {
	binary := strings.Replace(module.Path, ".go", ".so", 1)
	p, err := plugin.Load(binary)
	if err != nil {
		return err
	}
	mod := p.(plugin.ModuleInterface)
	if mod == nil {
		return fmt.Errorf("cannot read module")
	}
	err = mod.Connect(m.ctrl)
	if err != nil {
		return err
	}
	setup, err := mod.Setup()
	if err != nil {
		return err
	}

	module.Name = setup.Name
	module.Type = setup.Type
	module.Version = setup.Version
	module.Author = setup.Author
	module.Description = setup.Description

	m.modules[module.Name] = mod

	log.Event("Module '%s' loaded.", module.Name)
	return nil
}
