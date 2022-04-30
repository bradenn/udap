// Copyright (c) 2021 Braden Nicholson

package models

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"os/exec"
	"strings"
	"time"
	"udap/internal/store"
)

type Module struct {
	store.Persistent
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
}

func (m *Module) Build() error {
	// Create output file by modifying input file extension
	out := strings.Replace(m.Path, ".go", ".so", 1)
	// Create a timeout to prevent modules from taking too long to build
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	// Cancel the timeout of it exits before the timeout is up
	defer cancelFunc()
	// Prepare the command arguments
	args := []string{"build", "-v", "-buildmode=plugin", "-o", out, m.Path}
	// Initialize the command structure
	cmd := exec.CommandContext(timeout, "go", args...)
	// Run and get the stdout and stderr from the output
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (m *Module) Load() error {

	return nil
}

func (m *Module) Run() error {

	return nil
}

func (m *Module) Unload() error {
	return nil
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
func (m *Module) Emplace() (err error) {
	m.UpdatedAt = time.Now()
	err = store.DB.Model(&Module{}).Where("path = ?", m.Path).FirstOrCreate(&m).Error
	if err != nil {
		return err
	}
	return nil
}
