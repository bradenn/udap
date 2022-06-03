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
	Frequency  time.Duration
	*controller.Controller
}

// UpdateInterval is called once at the launch of the module
func (m *Module) UpdateInterval(frequency time.Duration) error {
	m.LastUpdate = time.Now()
	m.Frequency = time.Millisecond * frequency
	return nil
}

// Ready is called once at the launch of the module
func (m *Module) Ready() bool {
	return time.Since(m.LastUpdate) >= m.Frequency
}

// Connect is called once at the launch of the module
func (m *Module) Connect(ctrl *controller.Controller) error {
	m.LastUpdate = time.Now()
	m.Controller = ctrl
	return nil
}
