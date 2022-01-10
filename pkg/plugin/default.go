// Copyright (c) 2021 Braden Nicholson

package plugin

import "udap/internal/bond"

type Config struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // Module, Daemon, etc.
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
}

type Module struct {
	Config
	*bond.Bond
}

// Connect is called once at the launch of the module
func (m *Module) Connect(b *bond.Bond) error {
	m.Bond = b
	return nil
}
