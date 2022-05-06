// Copyright (c) 2022 Braden Nicholson

package server

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"udap/internal/bond"
	"udap/internal/controller"
	"udap/internal/log"
	"udap/internal/pulse"
	"udap/pkg/plugin"
)

const DIR = "modules"

type Modules struct {
	modules    map[string]plugin.ModuleInterface
	ctrl       *controller.Controller
	bond       *bond.Bond
	configured bool
	running    bool
}

func (m *Modules) Name() string {
	return "modules"
}

func (m *Modules) buildModules() error {
	// Try to load modules from the plugin folders
	err := m.buildModuleDir(DIR)
	if err != nil {
		return err
	}
	return nil
}

func (m *Modules) values() (pls []plugin.ModuleInterface, err error) {
	for _, up := range m.modules {
		pls = append(pls, up)
	}
	return pls, nil
}

func (m *Modules) Setup(ctrl *controller.Controller, bond *bond.Bond) error {
	m.ctrl = ctrl
	m.bond = bond
	m.configured = false
	m.modules = map[string]plugin.ModuleInterface{}
	err := m.buildModules()
	if err != nil {
		return err
	}
	return nil
}

func (m *Modules) Run() error {
	// Attempt to load the modules in the directory 'modules'
	err := m.loadModulesDir(DIR)
	if err != nil {
		return err
	}

	values, err := m.values()
	if err != nil {
		return err
	}

	// Create a wait group so all plugins can init at the same time
	wg := sync.WaitGroup{}
	wg.Add(len(values))
	// Run the full lifecycle of all plugins
	for _, module := range values {
		// Run a go function to create a new thread
		go func(p plugin.ModuleInterface) {
			defer wg.Done()
			if p == nil {
				log.Err(fmt.Errorf("invalid plugin"))
				return
			}
			// Defer the wait group to complete at the end
			// Attempt to connect to the module
			err = p.Connect(m.ctrl)
			if err != nil {
				return
			}
			// Run module setup
			c := plugin.Config{}
			c, err = p.Setup()
			if err != nil {
				log.ErrF(err, "Module '%s' setup failed: ", c.Name)
				return
			}
			start := time.Now()
			// Attempt to run the module
			err = p.Run()
			if err != nil {
				log.ErrF(err, "Module '%s' terminated prematurely: ", c.Name)
				return
			}
			log.Event("Module '%s' loaded. (%s)", c.Name, time.Since(start))
		}(module)
	}
	wg.Wait()

	return nil
}

func (m *Modules) Update() error {
	pulse.Fixed(2000)
	defer pulse.End()
	values, err := m.values()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(values))
	for _, module := range values {
		go func(k plugin.ModuleInterface) {
			defer wg.Done()
			err := k.Update()
			if err != nil {
				log.Err(err)
			}
		}(module)
	}
	wg.Wait()

	return nil
}

// buildModuleDir builds all potential modules in a directory
func (m *Modules) buildModuleDir(dir string) error {
	// Format the pattern for glob search
	pattern := fmt.Sprintf("./%s/*/*.go", dir)
	// Run the search for go files
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	// Add all the potential files from the search
	wg := &sync.WaitGroup{}
	wg.Add(len(files))
	// Launch a go func to build each one
	for _, p := range files {
		// Run the function for this file
		go func(path string) {
			defer wg.Done()
			if err = m.buildFromSource(path); err != nil {
				// If an error occurs, print it to console
				log.ErrF(err, "compiling module at '%s' failed", path)
				return
			}
			log.Event("Module '%s' compiled.", filepath.Base(path))
		}(p)
	}
	wg.Wait()
	return nil
}

// buildFromSource will build an eligible plugin from sources if applicable
func (m *Modules) buildFromSource(path string) error {
	// Create output file by modifying input file extension
	out := strings.Replace(path, ".go", ".so", 1)
	// Create a timeout to prevent modules from taking too long to build
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
	// Cancel the timeout of it exits before the timeout is up
	defer cancelFunc()
	// Prepare the command arguments
	args := []string{"build", "-v", "-buildmode=plugin", "-o", out, path}
	// Initialize the command structure
	cmd := exec.CommandContext(timeout, "go", args...)
	// Run and get the stdout and stderr from the output
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// loadModulesDir attempts to load each module
func (m *Modules) loadModulesDir(dir string) error {
	path := fmt.Sprintf("./%s/*/*.so", dir)
	files, err := filepath.Glob(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		p, err := plugin.Load(file)
		if err != nil {
			log.Err(err)
			continue
		}
		name := strings.Replace(filepath.Base(file), ".so", "", 1)
		mod := p.(plugin.ModuleInterface)

		m.modules[name] = mod

	}
	return nil
}
