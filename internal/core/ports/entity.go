// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type EntityRepository interface {
	common.Persist[domain.Entity]
	Register(*domain.Entity) error
	FindByName(name string) (*domain.Entity, error)
	FindAllByModule(name string) (*[]domain.Entity, error)
	FindAll() (*[]domain.Entity, error)
}

type EntityService interface {
	domain.Observable
	FindAll() (*[]domain.Entity, error)
	FindById(id string) (*domain.Entity, error)
	FindByName(name string) (*domain.Entity, error)
	FindAllByModule(name string) (*[]domain.Entity, error)
	Create(*domain.Entity) error
	ChangeIcon(id string, icon string) error
	SetPrediction(id string, prediction string) error
	ChangeAlias(id string, icon string) error
	Config(id string, value string) error
	FindOrCreate(*domain.Entity) error
	Register(*domain.Entity) error
	Update(*domain.Entity) error
	Delete(*domain.Entity) error
}
