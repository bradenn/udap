// Copyright (c) 2021 Braden Nicholson

package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
	"udap/internal/store"
)

type Module struct {
	store.Persistent
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
	Author      string      `json:"author"`
	Channel     chan Module `json:"-" gorm:"-"`
	State       string      `json:"state"`
	Enabled     bool        `json:"enabled" gorm:"default:true"`
}

// create inserts the current module into the database
func (m *Module) Register() error {
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
func (m *Module) Update() (err error) {
	m.UpdatedAt = time.Now()
	err = store.DB.Model(&Module{}).Where("id = ?", m.Id).Save(&m).Error
	if err != nil {
		return err
	}
	if m.Channel != nil {
		m.Channel <- *m
	}

	return nil
}

// Emplace gets a module from its path
func (m *Module) Emplace() (err error) {
	m.UpdatedAt = time.Now()
	err = store.DB.Model(&Module{}).Where("path = ?", m.Path).FirstOrCreate(&m).Error
	if err != nil {
		return err
	}
	if m.Channel != nil {
		m.Channel <- *m
	}
	return nil
}
