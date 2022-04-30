// Copyright (c) 2021 Braden Nicholson

package plugin

import (
	"time"
	"udap/internal/controller"
)

type Config struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // Module, Daemon, etc.
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
}

type Module struct {
	Config
	LastUpdate time.Time
	Frequency  int
	*controller.Controller
}

// Connect is called once at the launch of the module
func (m *Module) Connect(ctrl *controller.Controller) error {
	m.LastUpdate = time.Now()
	m.Controller = ctrl
	return nil
}
