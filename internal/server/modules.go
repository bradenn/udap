// Copyright (c) 2022 Braden Nicholson

package server

import (
	"context"
	"fmt"
	"os"
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

type Modules struct {
	modules map[string]plugin.UdapPlugin
	ctrl    *controller.Controller
	bond    *bond.Bond
	running bool
}

func (m *Modules) Name() string {
	return "modules"
}

func (m *Modules) buildModules() error {
	// Try to load modules from the plugin folders
	err := m.buildModuleDir("modules")
	if err != nil {
		return err
	}
	return nil
}

func (m *Modules) values() (pls []plugin.UdapPlugin, err error) {
	for _, up := range m.modules {
		pls = append(pls, up)
	}
	return pls, nil
}

func (m *Modules) Setup(ctrl *controller.Controller, bond *bond.Bond) error {
	m.ctrl = ctrl
	m.bond = bond
	m.modules = map[string]plugin.UdapPlugin{}
	err := m.buildModules()
	if err != nil {
		return err
	}
	return nil
}

func (m *Modules) Run() error {
	// Attempt to load the modules in the directory 'modules'
	err := m.loadModulesDir("modules")
	if err != nil {
		return err
	}
	// Create a wait group so all plugins can init at the same time
	wg := sync.WaitGroup{}
	values, err := m.values()
	if err != nil {
		return err
	}
	wg.Add(len(values))
	// Run the full lifecycle of all plugins
	for _, module := range values {
		// Run a go function to create a new thread
		go func(p plugin.UdapPlugin) {
			defer wg.Done()
			if p == nil {
				log.Err(fmt.Errorf("invalid plugin"))
				return
			}
			// Defer the wait group to complete at the end
			// Attempt to connect to the module
			err = p.Connect(m.ctrl, m.bond)
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
			log.Event("Module '%s' v%s running.", c.Name, c.Version)
			// Attempt to run the module
			err = p.Run()
			if err != nil {
				log.ErrF(err, "Module '%s' terminated prematurely: ", c.Name)
				return
			}
			log.Event("Module '%s' exited. (%s)", c.Name, time.Since(start))
		}(module)
	}
	wg.Wait()
	return nil
}

func (m *Modules) Update() error {
	pulse.Fixed(500)
	defer pulse.End()
	values, err := m.values()
	if err != nil {
		return err
	}

	for _, mod := range values {
		err = mod.Update()
		if err != nil {
			return err
		}
	}
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
				log.ErrF(err, "failed to build module candidate '%s'", path)
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
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	// Cancel the timeout of it exits before the timeout is up
	defer cancelFunc()
	// get the go executable from the environment
	goExec := os.Getenv("goExec")
	// Prepare the command arguments
	args := []string{"build", "-v", "-buildmode=plugin", "-o", out, path}
	// Initialize the command structure
	cmd := exec.CommandContext(timeout, goExec, args...)
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
		mod := p.(plugin.UdapPlugin)
		m.modules[name] = mod

	}
	return nil
}
