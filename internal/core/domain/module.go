// Copyright (c) 2022 Braden Nicholson

package domain

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
	Enabled     bool        `json:"enabled" gorm:"default:true"`
	Recover     int         `json:"recover"`
}

type ModuleRepository interface {
	FindAll() (*[]Module, error)
	FindByName(name string) (*Module, error)
	FindById(id string) (*Module, error)
	Create(*Module) error
	FindOrCreate(*Module) error
	Update(*Module) error
	Delete(*Module) error
}

type ModuleOperator interface {
	Build(module *Module) error
	Load(module *Module) error
	Update(module *Module) error
	Run(module *Module) error
}

type ModuleService interface {
	Discover() error
	Build(module *Module) error
	Load(module *Module) error
	Update(module *Module) error
	Run(module *Module) error
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
