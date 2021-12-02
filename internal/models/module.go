// Copyright (c) 2021 Braden Nicholson

package models

import (
	"fmt"
	"gorm.io/gorm"
	"udap/internal/store"
	"udap/pkg/plugin"
)

type Module struct {
	store.Persistent
	plugin.Plugin   `gorm:"-" json:"-"`
	plugin.Metadata `gorm:"-"`
	// Path refers to the literal name of the module
	Path string `json:"-"`
}

func (m *Module) Event(event plugin.Event) (err error) {
	return nil
}

// FromPath gets a module from its path
func FromPath(path string, p plugin.Plugin) (module Module, err error) {
	module.Path = path
	err = module.setPlugin(p)
	if err != nil {
		return Module{}, err
	}
	err = store.DB.Model(&Module{}).Where("path = ?", path).FirstOrCreate(&module).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return module, err
		}
	}
	return module, nil
}

// SetPlugin sets the module's reference to it's plugin
func (m *Module) setPlugin(p plugin.Plugin) (err error) {
	// Make sure the plugin isn't already set
	if m.Plugin != nil {
		return fmt.Errorf("module plugin is already set")
	}
	//
	m.Plugin = p
	m.Metadata = p.Metadata()
	return nil
}

// create inserts the current module into the database
func (m *Module) registerEntity(name string) error {
	// Attempt to create a new module
	err := store.DB.Model(&Entity{}).Where("name = ? AND moduleId = ?", name, m.Id).FirstOrCreate(m).Error
	// Report internal errors for later diagnostic
	if err != nil {
		return fmt.Errorf("failed to create module")
	}
	// Return no errors
	return nil
}

// create inserts the current module into the database
func (m *Module) create() error {
	// Attempt to create a new module
	err := store.DB.Create(m).Error
	// Report internal errors for later diagnostic
	if err != nil {
		return fmt.Errorf("failed to create module")
	}
	// Return no errors
	return nil
}

// Hooks

func (m *Module) BeforeCreate(_ *gorm.DB) error {
	if m.Plugin == nil {
		return fmt.Errorf("plugin does not exist; cannot create module")
	}
	m.Metadata = m.Plugin.Metadata()
	return nil
}

func (m *Module) AfterFind(_ *gorm.DB) error {

	return nil
}
