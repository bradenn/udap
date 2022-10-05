// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type ModuleRepository interface {
	common.Persist[domain.Module]
	FindByName(name string) (*domain.Module, error)
	FindByUUID(uuid string) (*domain.Module, error)
}

type ModuleOperator interface {
	Build(module string, uuid string) error
	Load(module string, uuid string) (domain.ModuleConfig, error)
	Dispose(module string, uuid string) error
	Run(uuid string) error
	Update(uuid string) error
	HandleEmit(mutation domain.Mutation) error
}

type ModuleService interface {
	domain.Observable
	Discover() error
	InitConfig(string, string, string) error
	SetConfig(string, string, string) error
	GetConfig(string, string) (string, error)
	HandleEmits(mutation domain.Mutation) error
	Build(id string) error
	Load(id string) error
	Update(id string) error
	Run(id string) error
	Dispose(id string) error
	UpdateAll() error
	RunAll() error
	DisposeAll() error
	LoadAll() error
	BuildAll() error
	FindAll() (*[]domain.Module, error)
	FindByName(name string) (*domain.Module, error)
	Disable(name string) error
	Enable(name string) error
	Reload(name string) error
	Halt(name string) error
}
