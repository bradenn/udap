// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type ZoneRepository interface {
	common.Persist[domain.Zone]
	FindByName(name string) (*domain.Zone, error)
}

type ZoneService interface {
	domain.Observable
	AddEntity(id string, entity string) error
	RemoveEntity(id string, entity string) error
	Pin(id string) error
	Unpin(id string) error
	FindAll() (*[]domain.Zone, error)
	FindById(id string) (*domain.Zone, error)
	FindByName(name string) (*domain.Zone, error)
	Create(*domain.Zone) error
	FindOrCreate(*domain.Zone) error
	Update(*domain.Zone) error
	Delete(id string) error
	Restore(id string) error
}
