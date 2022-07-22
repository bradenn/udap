// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type Zone struct {
	common.Persistent
	Name     string   `json:"name"`
	Entities []Entity `json:"entities" gorm:"many2many:zone_entities;"`
	User     string   `json:"user"`
}

type ZoneRepository interface {
	common.Persist[Zone]
	FindByName(name string) (*Zone, error)
}

type ZoneService interface {
	Observable
	FindAll() (*[]Zone, error)
	FindById(id string) (*Zone, error)
	FindByName(name string) (*Zone, error)
	Create(*Zone) error
	FindOrCreate(*Zone) error
	Update(*Zone) error
	Delete(id string) error
	Restore(id string) error
}
