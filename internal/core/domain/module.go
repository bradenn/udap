// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/pkg/plugin"

type Module struct {
	Persistent
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
	Author      string      `json:"author"`
	Channel     chan Module `json:"-" gorm:"-"`
	State       string      `json:"state"`
	plugin.ModuleInterface
	Enabled bool `json:"enabled" gorm:"default:true"`
	Recover int  `json:"recover"`
}

type ModuleRepository interface {
	Candidates() ([]string, error)
	FindAll() ([]*Module, error)
	FindByName(name string) (*Module, error)
}

type ModuleService interface {
	Discover() error
	Build(name string) error
	BuildAll() error
	FindAll() ([]*Module, error)
	FindByName(name string) (*Module, error)
	Disable(name string) error
	Enable(name string) error
	Reload(name string) error
	Halt(name string) error
}
