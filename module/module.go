// Copyright (c) 2021 Braden Nicholson

package module

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"plugin"
	"strings"
	"udap/udap"
	"udap/udap/db"
	"udap/udap/store"
)

var modules map[string]*Module

func init() {
	modules = map[string]*Module{}
}

func Get(id string) (mod *Module, err error) {
	mod = modules[id]
	if mod == nil {

		return nil, fmt.Errorf("module is nil")
	}
	return mod, nil
}

const (
	pathFmt   = "./plugins"
	pluginFmt = pathFmt + "/%s/%s.so"
)

func (m Module) path() string {
	return fmt.Sprintf("modules.%s", strings.ToLower(m.Name))
}

type Replicant struct {
	config map[string]string
}

type Plugin interface {
	// Startup is called when the plugin is first loaded
	Startup() (Metadata, error)
	// Default returns the plugin's default instance config
	Default() interface{}

	Update(ctx Context) error
	Run(ctx Context, data string) (string, error)
}

type Context struct {
	update     chan UpdateBuffer
	instanceId string
}

func (c *Context) Send(data interface{}) {
	marshaled, err := json.Marshal(data)
	if err != nil {
		return
	}
	c.update <- UpdateBuffer{
		InstanceId: c.instanceId,
		Data:       string(marshaled),
	}
}

func (c *Context) Id() string {
	return c.instanceId
}

type Module struct {
	udap.Persistent
	// Metadata holds generic information about the module.
	Metadata
	// Environment holds persistent JSON environment information, api keys, ports, etc.
	Environment string `json:"environment"`
	// Path refers to the literal name of the module
	Path string `json:"path"`

	component Plugin `gorm:"-"`
}

// Emplace interface operations
func (m *Module) Emplace() error {

	component, err := m.getComponent()
	if err != nil {
		return err
	}
	// Populate the plugin's metadata field
	m.component = component
	metadata, err := m.component.Startup()
	if err != nil {
		return err
	}
	m.Metadata = metadata
	marshal, err := json.Marshal(m.Metadata)
	if err != nil {
		return err
	}

	// Find the module if it exists in the database
	err = db.DB.Model(&Module{}).Where("path = ?", m.Path).First(&m).Error
	if err != nil {
		// If the module does not exist, create it
		if err == gorm.ErrRecordNotFound {
			err = m.create()
			// Report errors regarding the creation
			if err != nil {
				return err
			}

		}
		// Return any unhandled errors to be reported
		return err
	}
	modules[m.Id] = m
	err = store.PutLn(string(marshal), m.path(), "metadata")
	if err != nil {
		return err
	}
	// If the module has been loaded, or created announce it
	udap.Info("Loaded module '%s' v%s", m.Name, m.Version)
	return nil
}

// Emplace interface operations
func (m *Module) Run(data string) error {

	_, err := m.component.Run(Context{}, data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// create inserts the current module into the database
func (m *Module) create() error {
	err := db.DB.Create(m).Error
	// Report internal errors for later diagnostic
	if err != nil {
		return fmt.Errorf("failed to create module")
	}
	return nil
}

func (m *Module) getComponent() (Plugin, error) {

	filename := fmt.Sprintf("./plugins/%s", m.Path)
	p, err := plugin.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("invalid or out-dated plugin '%s'", m.Path)
	}

	lookup, err := p.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("plugin '%s' provides not exported accesspoint", m.Path)
	}
	return lookup.(Plugin), nil
}

func (m *Module) valid() bool {
	component, _ := m.getComponent()
	if component == nil {
		return false
	}
	return true
}

// Hooks

func (m *Module) BeforeCreate(_ *gorm.DB) error {

	return nil
}

func (m *Module) AfterFind(_ *gorm.DB) error {

	return nil
}
