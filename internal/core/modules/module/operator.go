// Copyright (c) 2022 Braden Nicholson

package module

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"strings"
	"time"
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/plugin"
	"udap/internal/pulse"
)

type moduleOperator struct {
	ctrl    *controller.Controller
	runtime map[string]plugin.ModuleInterface
}

func NewOperator(ctrl *controller.Controller) domain.ModuleOperator {
	return &moduleOperator{
		ctrl:    ctrl,
		runtime: map[string]plugin.ModuleInterface{},
	}
}

func (m *moduleOperator) getModule(id string) (plugin.ModuleInterface, error) {
	if m.runtime[id] == nil {
		return nil, fmt.Errorf("module not found")
	}
	return m.runtime[id], nil
}

func (m *moduleOperator) setModule(id string, moduleInterface plugin.ModuleInterface) error {
	m.runtime[id] = moduleInterface
	return nil
}

func (m *moduleOperator) removeModule(id string) error {
	delete(m.runtime, id)
	return nil
}

func (m *moduleOperator) alreadyExists(oldPath string, newPath string) error {
	oldStats, err := os.Stat(oldPath)
	if err != nil {
		return err
	}
	newStats, err := os.Stat(newPath)
	if err != nil {
		return err
	}
	if oldStats.Size() == newStats.Size() {
		return fmt.Errorf("SAME FILE")
	}
	return nil
}

// Build will compile a valid plugin file into a readable binary
func (m *moduleOperator) Build(module *domain.Module) error {
	start := time.Now()
	if _, err := os.Stat(module.Path); err != nil {
		return err
	}
	// Create a timeout to prevent modules from taking too long to build
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
	// Cancel the timeout of it exits before the timeout is up
	defer cancelFunc()

	module.UUID = uuid.New().String()
	// Create the binary file path
	binary := strings.Replace(module.Path, ".go", fmt.Sprintf("-%s.so", module.UUID), 1)
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

	log.Event("Module '%s' @ %s compiled. (%s)", module.Name, module.SessionId(),
		time.Since(start).Truncate(time.Millisecond).
			String())
	return nil
}

// Load is used to find a pre-built plugin file, and load it into the local system.
// The module reference should be saved to the repository after loading.
func (m *moduleOperator) Load(module *domain.Module) error {
	// Create the binary file path
	binary := module.CompiledPath()
	// Attempt to load the plugin binary
	p, err := plugin.Load(binary)
	if err != nil {
		return err
	}
	// Extract the plugin interface
	mod := p.(plugin.ModuleInterface)
	if mod == nil {
		return fmt.Errorf("cannot read module")
	}
	// Connect the module to the UDAP runtime
	err = mod.Connect(m.ctrl)
	if err != nil {
		return err
	}
	// Run the setup method
	setup, err := mod.Setup()
	if err != nil {
		return err
	}
	// Update the module parameters
	module.Name = setup.Name
	module.Type = setup.Type
	module.Version = setup.Version
	module.Author = setup.Author
	module.Description = setup.Description
	// Emplace the module into the local buffer
	err = m.setModule(module.Id, mod)
	if err != nil {
		return err
	}
	// Log the status
	log.Event("Module '%s' @ %s loaded.", module.Name, module.SessionId())
	return nil
}

// Run attempts to initialize the plugin's runtime
func (m *moduleOperator) Run(module *domain.Module) error {
	// Check to make sure the module is enabled
	if !module.Enabled {
		return fmt.Errorf("module '%s' is not enabled; cannot run", module.Name)
	}
	// Make sure the module is not already running
	if module.Running {
		return fmt.Errorf("module '%s' is already running; cannot run again", module.Name)
	}
	// Get the local module
	local, err := m.getModule(module.Id)
	if err != nil {
		return err
	}
	// Run the local module
	err = local.Run()
	if err != nil {
		return err
	}

	return nil
}

// Update is called on every tick, it is the plugin's decision whether to use the time or defer.
func (m *moduleOperator) Update(module *domain.Module) error {
	// Check to make sure the module is enabled
	if !module.Enabled {
		return fmt.Errorf("module '%s' is not enabled; cannot update", module.Name)
	}
	// Make sure the module is running
	if !module.Running {
		return fmt.Errorf("module '%s' is not running; cannot update", module.Name)
	}
	// Get the local module
	local, err := m.getModule(module.Id)
	if err != nil {
		return err
	}
	// Begin the lifecycle metrics
	pulse.Begin(module.Id)
	// Run the update
	err = local.Update()
	if err != nil {
		return err
	}
	// End the pulse metric
	pulse.End(module.Id)

	return nil
}

// Dispose is called at the end of the lifecycle, it attempts to halt activity.
func (m *moduleOperator) Dispose(module *domain.Module) error {
	// Check to make sure the module is enabled
	if !module.Enabled {
		return fmt.Errorf("module '%s' is not enabled; cannot dispose", module.Name)
	}
	// Make sure the module is running
	if !module.Running {
		return fmt.Errorf("module '%s' is not running; cannot dispose", module.Name)
	}
	// Get the local module
	local, err := m.getModule(module.Id)
	if err != nil {
		return err
	}
	// Dispose of the local module
	err = local.Dispose()
	if err != nil {
		return err
	}
	defer func() {
		err = os.Remove(module.CompiledPath())
		if err != nil {
			return
		}
	}()
	err = m.removeModule(module.Id)
	if err != nil {
		return err
	}
	log.Event("Module '%s' @ %s unloaded.", module.Name, module.SessionId())
	module.UUID = ""
	return nil
}
