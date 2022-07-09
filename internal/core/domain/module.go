// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type Module struct {
	common.Persistent
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
	Author      string      `json:"author"`
	Channel     chan Module `json:"-" gorm:"-"`
	State       string      `json:"state"`
	Enabled     bool        `json:"enabled" gorm:"default:true"`
	Recover     int         `json:"recover"`
}

type ModuleRepository interface {
	common.Persist[Module]
	FindByName(name string) (*Module, error)
}

type ModuleOperator interface {
	Build(module *Module) error
	Load(module *Module) error
	Update(module *Module) error
	Run(module *Module) error
	Dispose(module *Module) error
}

type ModuleService interface {
	Observable
	Discover() error
	Build(module *Module) error
	Load(module *Module) error
	Update(module *Module) error
	Run(module *Module) error
	Dispose(module *Module) error
	UpdateAll() error
	RunAll() error
	LoadAll() error
	BuildAll() error
	FindAll() (*[]Module, error)
	FindByName(name string) (*Module, error)
	Disable(name string) error
	Enable(name string) error
	Reload(name string) error
	Halt(name string) error
}
