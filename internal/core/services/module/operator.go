// Copyright (c) 2022 Braden Nicholson

package module

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/plugin"
)

type moduleOperator struct {
	ctrl    *controller.Controller
	runtime map[string]plugin.ModuleInterface
}

func (m *moduleOperator) HandleEmit(mutation domain.Mutation) error {
	return nil
}

const PATH = "./modules"

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

// hashFile returns a sha256 hash of a source file
func hashFile(source string) (string, error) {
	f, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}
	return string(h.Sum(nil)), nil
}

// generateSessionId provides the first bytes of the module uuid
func generateSessionId(uuid string) string {
	return strings.Split(uuid, "-")[0]
}

// generateSourcePath generates the path to a module's source code
func generateSourcePath(module string) string {
	return fmt.Sprintf("%s/%s/%s.go", PATH, module, module)
}

// generateBuildPath generates the path to a module's binary file
func generateBuildPath(module string, uuid string) string {
	return fmt.Sprintf("%s/%s/%s-%s.so", PATH, module, module, uuid)
}

// Build will compile a valid plugin file into a readable binary
func (m *moduleOperator) Build(module string, uuid string) error {
	// Generate the source file path from the module name
	sourcePath := generateSourcePath(module)
	// Generate the output build file path
	buildPath := generateBuildPath(module, uuid)
	// Confirm that the source file exists
	if _, err := os.Stat(sourcePath); err != nil {
		return err
	}
	// Create a timeout to prevent modules from taking too long to build
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	// Cancel the timeout of it exits before the timeout is up
	defer cancelFunc()
	// Prepare the command arguments
	args := []string{"build", "-v", "-buildmode=plugin", "-o", buildPath, sourcePath}
	// Initialize the command structure
	cmd := exec.CommandContext(timeout, "go", args...)
	// Run and get the stdout and stderr from the output
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("module build failed: \n%s", out)
	}
	return nil
}

// Load is used to find a pre-built plugin file, and load it into the local system.
// The module reference should be saved to the repository after loading.
func (m *moduleOperator) Load(module string, uuid string) (domain.ModuleConfig, error) {
	// Create the binary file path
	binary := generateBuildPath(module, uuid)
	// Attempt to load the plugin binary
	p, err := plugin.Load(binary)
	if err != nil {
		return domain.ModuleConfig{}, err
	}
	// Extract the plugin interface
	mod := p.(plugin.ModuleInterface)
	if mod == nil {
		return domain.ModuleConfig{}, fmt.Errorf("cannot read module")
	}
	// Connect the module to the UDAP runtime
	err = mod.Connect(m.ctrl)
	if err != nil {
		return domain.ModuleConfig{}, err
	}
	// Run the setup method
	setup, err := mod.Setup()
	if err != nil {
		return domain.ModuleConfig{}, err
	}
	// Emplace the module into the local buffer
	err = m.setModule(uuid, mod)
	if err != nil {
		return domain.ModuleConfig{}, err
	}
	conf := domain.ModuleConfig{
		Name:        setup.Name,
		Type:        setup.Type,
		Description: setup.Description,
		Version:     setup.Version,
		Author:      setup.Author,
	}
	return conf, nil
}

// Dispose is called at the end of the lifecycle, it attempts to halt activity.
func (m *moduleOperator) Dispose(module string, uuid string) error {
	// Get the local module
	local, err := m.getModule(uuid)
	if err != nil {
		return err
	}
	binaryPath := generateBuildPath(module, uuid)
	// Confirm that the binary file exists
	if _, err = os.Stat(binaryPath); err != nil {
		return err
	}
	// Dispose of the local module
	err = local.Dispose()
	if err != nil {
		return err
	}
	// Remove the file when the function exits
	defer func() {
		err = os.Remove(binaryPath)
		if err != nil {
			log.Err(err)
			return
		}
	}()
	err = m.removeModule(uuid)
	if err != nil {
		return err
	}
	return nil
}

// Run attempts to initialize the plugin's runtime
func (m *moduleOperator) Run(uuid string) error {
	// Get the local module
	local, err := m.getModule(uuid)
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
func (m *moduleOperator) Update(uuid string) error {
	// Get the local module
	local, err := m.getModule(uuid)
	if err != nil {
		return err
	}
	// Run the update
	err = local.Update()
	if err != nil {
		return err
	}
	return nil
}
