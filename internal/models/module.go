// Copyright (c) 2021 Braden Nicholson

package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
	"udap/internal/store"
	"udap/pkg/plugin"
)

type Module struct {
	store.Persistent
	plugin.Module
	plugin.UdapPlugin `json:"-" gorm:"-"`
	// Path refers to the literal name of the module
	Path string `json:"-"`
}

func (m *Module) Event(event plugin.Event) (err error) {
	return nil
}

// FromPath gets a module from its path
func FromPath(path string, p plugin.Plugin) (module Module, err error) {
	module.Path = path
	// err = module.setPlugin(p)
	// if err != nil {
	// 	return Module{}, err
	// }
	err = store.DB.Model(&Module{}).Where("path = ?", path).FirstOrCreate(&module).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return module, err
		}
	}
	return module, nil
}

// create inserts the current module into the database
func (m *Module) registerEntity(entity *Entity) error {
	err := entity.Find()
	if err != nil {
		return err
	}

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

	return nil
}

func (m *Module) AfterFind(_ *gorm.DB) error {

	return nil
}

// Emplace gets a module from its path
func (m *Module) Emplace() (err error) {
	m.UpdatedAt = time.Now()
	err = store.DB.Model(&Module{}).Where("path = ?", m.Path).FirstOrCreate(&m).Error
	if err != nil {
		return err
	}
	return nil
}
