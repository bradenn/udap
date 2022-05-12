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
	"udap/internal/models"
	"udap/internal/pulse"
	"udap/pkg/plugin"
)

const DIR = "modules"

const (
	UNINITIALIZED = "uninitialized"
	IDLE          = "idle"
	RUNNING       = "running"
	STOPPED       = "stopped"
	ERROR         = "error"
)

type ModuleController struct {
	model   models.Module
	loaded  bool
	running bool
	module  plugin.ModuleInterface
	config  plugin.Config
	source  string
	binary  string
	state   string
	c       chan ModuleState
	mid     string

	receive chan models.Module
}

func (m *ModuleController) listen() {
	for module := range m.receive {
		if m.model.Enabled != module.Enabled {
			if !module.Enabled {
				log.Event("Disabling module '%s'", m.model.Name)
			} else {
				log.Event("Enabling module '%s'", m.model.Name)
			}
		}
		m.model = module
	}
}

func NewModuleController(path string) (*ModuleController, error) {
	if !strings.Contains(path, ".go") {
		return nil, fmt.Errorf("invald module source path provided")
	}
	recv := make(chan models.Module)
	m := &ModuleController{
		running: false,
		loaded:  false,
		module:  nil,
		config:  plugin.Config{},
		source:  path,
		receive: recv,
		binary:  strings.Replace(path, ".go", ".so", 1),
	}
	go m.listen()
	return m, nil
}

func (m *ModuleController) setState(state string) {
	m.state = state
	if m.c == nil {
		return
	}
	m.c <- ModuleState{
		id:    m.mid,
		state: m.state,
	}
}

func (m *ModuleController) build() error {
	// Create a timeout to prevent modules from taking too long to build
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
	// Cancel the timeout of it exits before the timeout is up
	defer cancelFunc()
	// Prepare the command arguments
	args := []string{"build", "-v", "-buildmode=plugin", "-o", m.binary, m.source}
	// Initialize the command structure
	cmd := exec.CommandContext(timeout, "go", args...)
	// Run and get the stdout and stderr from the output
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (m *ModuleController) setup(ctrl *controller.Controller, c chan ModuleState) error {
	err := m.build()
	if err != nil {
		return err
	}
	m.state = UNINITIALIZED
	p, err := plugin.Load(m.binary)
	if err != nil {
		return err
	}
	mod := p.(plugin.ModuleInterface)
	if mod == nil {
		return fmt.Errorf("cannot read module")
	}
	m.module = mod
	// Defer the wait group to complete at the end
	// Attempt to connect to the module
	err = m.module.Connect(ctrl)
	if err != nil {
		return err
	}
	// Run module setup
	m.config, err = p.Setup()
	if err != nil {
		log.ErrF(err, "Module '%s' setup failed: ", m.config.Name)
		return err
	}

	modDb := models.Module{
		Name:        m.config.Name,
		Path:        m.source,
		Type:        m.config.Type,
		Description: m.config.Description,
		Version:     m.config.Version,
		Author:      m.config.Author,
		State:       m.state,
		Channel:     m.receive,
	}
	m.c = c
	mid, err := ctrl.Modules.Register(modDb)
	if err != nil {
		return err
	}
	modDb.Id = mid
	m.model = modDb
	m.mid = mid

	m.setState(IDLE)
	m.loaded = true
	return nil
}

func (m *ModuleController) start() error {
	if !m.loaded {
		return fmt.Errorf("module is not loaded, setup must be called first")
	}
	// Attempt to run the module
	m.setState(RUNNING)
	go func() {
		for {
			err := m.module.Run()
			if err != nil {
				log.ErrF(err, "Module '%s' terminated prematurely: ", m.config.Name)
				m.setState(ERROR)
				break
			}
			break
		}
		m.setState(IDLE)
	}()

	log.Event("Module '%s' loaded.", m.config.Name)
	m.running = true

	return nil
}

func (m *ModuleController) update() error {
	if !m.running {
		return nil
	}
	if !m.loaded {
		return fmt.Errorf("module is not loaded, setup must be called first")
	}
	if m.model.Enabled {
		err := m.module.Update()
		if err != nil {
			return err
		}
	}
	return nil
}

type ModuleState struct {
	id    string
	state string
}

type Modules struct {
	mcs        map[string]*ModuleController
	modMutex   sync.Mutex
	ctrl       *controller.Controller
	bond       *bond.Bond
	state      chan ModuleState
	configured bool
	running    bool
}

func (m *Modules) Name() string {
	return "modules"
}

func (m *Modules) listen() {

	for {
		select {
		case state := <-m.state:
			err := m.ctrl.Modules.State(state.id, state.state)
			if err != nil {
				log.Err(err)
				continue
			}
		}
	}

}

func (m *Modules) buildModules() error {
	// Try to load modules from the plugin folders
	err := m.buildModuleDir(DIR)
	if err != nil {
		return err
	}
	return nil
}

func (m *Modules) values() (pls []*ModuleController, err error) {
	for _, up := range m.mcs {
		pls = append(pls, up)
	}
	return pls, nil
}

func (m *Modules) Setup(ctrl *controller.Controller, bond *bond.Bond) error {
	m.ctrl = ctrl
	m.bond = bond
	m.configured = false
	m.mcs = map[string]*ModuleController{}
	m.modMutex = sync.Mutex{}
	m.state = make(chan ModuleState, 20)
	go m.listen()
	err := m.buildModules()
	if err != nil {
		return err
	}
	return nil
}

func (m *Modules) Run() error {
	// Create a wait group so all plugins can init at the same time
	m.modMutex.Lock()
	wg := sync.WaitGroup{}
	wg.Add(len(m.mcs))
	// Run the full lifecycle of all plugins
	for _, module := range m.mcs {
		// Run a go function to create a new thread
		go func(p *ModuleController) {
			defer wg.Done()
			err := p.setup(m.ctrl, m.state)
			if err != nil {
				log.ErrF(err, "Module '%s' setup failed: ", p.config.Name)
				return
			}
			// Attempt to run the module
			err = p.start()
			if err != nil {
				return
			}
		}(module)
	}
	wg.Wait()
	m.modMutex.Unlock()

	return nil
}

func (m *Modules) Update() error {
	wg := sync.WaitGroup{}
	m.modMutex.Lock()
	wg.Add(len(m.mcs))
	for _, module := range m.mcs {
		go func(p *ModuleController) {
			defer wg.Done()
			if p.mid == "" {
				return
			}
			pulse.Begin(p.mid)
			err := p.update()
			if err != nil {
				return
			}
			pulse.End(p.mid)
		}(module)
	}
	wg.Wait()
	m.modMutex.Unlock()
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
	// Launch a go func to build each one
	for _, p := range files {
		mc, err := NewModuleController(p)
		if err != nil {
			return err
		}
		mc.config.Name = strings.Replace(filepath.Base(p), ".go", "", 1)
		m.modMutex.Lock()
		m.mcs[mc.config.Name] = mc
		m.modMutex.Unlock()
	}
	return nil
}
