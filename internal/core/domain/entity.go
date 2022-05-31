// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type Entity struct {
	common.Persistent
	Name      string `gorm:"unique" json:"name"` // Given name from module
	Alias     string `json:"alias"`              // Name from users
	Type      string `json:"type"`               // Type of entity {Light, Sensor, Etc}
	Module    string `json:"module"`             // Parent Module name
	Locked    bool   `json:"locked"`             // Is the Entity state locked?
	Config    string `json:"config"`
	Position  string `json:"position" gorm:"default:'{}'"`
	Icon      string `json:"icon" gorm:"default:'􀛮'"` // The icon to represent this entity
	Frequency int    `json:"frequency" gorm:"default:3000"`
	Neural    string `json:"neural" gorm:"default:'inactive'"` // Parent Module name
	Predicted string `gorm:"-" json:"predicted"`               // scalar
}

type EntityRepository interface {
	common.Persist[Entity]
	Register(*Entity) error
}

type EntityService interface {
	Observable
	FindAll() (*[]Entity, error)
	FindById(id string) (*Entity, error)
	Create(*Entity) error
	Config(id string, value string) error
	FindOrCreate(*Entity) error
	Register(*Entity) error
	Update(*Entity) error
	Delete(*Entity) error
}
