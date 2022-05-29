// Copyright (c) 2022 Braden Nicholson

package domain

type Zone struct {
	Persistent
	Name     string   `json:"name"`
	Entities []Entity `json:"entities" gorm:"many2many:zone_entities;"`
	User     string   `json:"user"`
}

type ZoneRepository interface {
	FindAll() ([]*Zone, error)
	FindById(id string) (*Zone, error)
	Create(*Zone) error
	FindOrCreate(*Zone) error
	Update(*Zone) error
	Delete(*Zone) error
}

type ZoneService interface {
	FindAll() ([]*Zone, error)
	FindById(id string) (*Zone, error)
	Create(*Zone) error
	FindOrCreate(*Zone) error
	Update(*Zone) error
	Delete(*Zone) error
}
