package types

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"plugin"
	"udap/config"
)

const (
	pathFmt   = "./plugins"
	pluginFmt = pathFmt + "/%s/%s.so"
)

type IModule interface {
	Load(agent Agent) error

	Create(agent Agent) error

	Run(data string) (error, string)
	// Metadata returns basic information about the module.
	Metadata() Metadata
}

type Module struct {
	Persistent
	// Metadata holds generic information about the module.
	Metadata
	// Environment holds persistent JSON environment information, api keys, ports, etc.
	Environment string `json:"environment"`
	// Path refers to the literal name of the module
	Path string `json:"path"`

	component IModule `gorm:"-"`
}

// Emplace interface operations
func (m *Module) Emplace() error {
	component, err := m.getComponent()
	if err != nil {
		return err
	}
	// Populate the plugin's metadata field
	m.component = component
	m.Metadata = component.Metadata()
	// Find the module if it exists in the database
	err = db.Model(&Module{}).Where("path = ?", m.Path).First(&m).Error
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
	// If the module has been loaded, or created announce it
	config.Info("Loaded module '%s' v%s", m.Name, m.Version)
	return nil
}

// Emplace interface operations
func (m *Module) Run(data string) error {

	err, state := m.component.Run(data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ent := Entity{}
	err = json.Unmarshal([]byte(data), &ent)
	if err != nil {
		return err
	}

	db.Model(&Entity{}).Where("name = ?", ent.Name).Update("state", state)
	return nil
}

// create inserts the current module into the database
func (m *Module) create() error {
	err := db.Create(m).Error
	// Report internal errors for later diagnostic
	if err != nil {
		return fmt.Errorf("failed to create module")
	}
	return nil
}

// Initialize loads the module under an instance
func (m *Module) Initialize(agent Agent) (component IModule, err error) {
	if m.component != nil {
		return m.component, nil
	}
	component, err = m.getComponent()
	if err != nil {
		return nil, err
	}
	err = component.Load(agent)
	if err != nil {
		return nil, err
	}
	return component, nil
}

func (m *Module) rawComponent() (IModule, error) {
	if m.component != nil {
		return m.component, nil
	}
	return m.getComponent()
}
func (m *Module) getComponent() (IModule, error) {

	filename := fmt.Sprintf("./plugins/%s", m.Path)
	p, err := plugin.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("invalid or out-dated plugin '%s'", m.Path)
	}

	lookup, err := p.Lookup("IModule")
	if err != nil {
		return nil, fmt.Errorf("plugin '%s' provides not exported accesspoint", m.Path)
	}
	return lookup.(IModule), nil
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
	if !m.valid() {
		return fmt.Errorf("module does not exist")
	}
	component, err := m.getComponent()
	if err != nil {
		return err
	}
	m.Metadata = component.Metadata()
	return nil
}

func (m *Module) AfterFind(_ *gorm.DB) error {

	return nil
}
