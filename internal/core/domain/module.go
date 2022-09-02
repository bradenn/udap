// Copyright (c) 2022 Braden Nicholson

package domain

import (
	"fmt"
	"strings"
	"udap/internal/core/domain/common"
)

type ModuleConfig struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // Module, Daemon, etc.
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
}

type Module struct {
	common.Persistent
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	UUID        string            `json:"uuid"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	Author      string            `json:"author"`
	Channel     chan Module       `json:"-" gorm:"-"`
	Config      map[string]string `json:"config"`
	State       string            `json:"state"`
	Running     bool              `json:"running" gorm:"default:false"`
	Enabled     bool              `json:"enabled" gorm:"default:true"`
	Recover     int               `json:"recover"`
}

func (m *Module) SessionId() string {
	if m.UUID == "" {
		return "invalid"
	}
	return strings.Split(m.UUID, "-")[0]
}

func (m *Module) CompiledPath() string {
	if m.UUID == "" {
		return "invalid"
	}
	return strings.Replace(m.Path, ".go", fmt.Sprintf("-%s.so", m.UUID), 1)
}

type ModuleRepository interface {
	common.Persist[Module]
	FindByName(name string) (*Module, error)
}

type ModuleOperator interface {
	Build(module string, uuid string) error
	Load(module string, uuid string) (ModuleConfig, error)
	Dispose(module string, uuid string) error
	Run(uuid string) error
	Update(uuid string) error
	HandleEmit(mutation Mutation) error
}

type ModuleService interface {
	Observable
	Discover() error
	HandleEmits(mutation Mutation) error
	Build(module *Module) error
	Load(module *Module) error
	Update(module *Module) error
	Run(module *Module) error
	Dispose(module *Module) error
	UpdateAll() error
	RunAll() error
	DisposeAll() error
	LoadAll() error
	BuildAll() error
	FindAll() (*[]Module, error)
	FindByName(name string) (*Module, error)
	Disable(name string) error
	Enable(name string) error
	Reload(name string) error
	Halt(name string) error
}
