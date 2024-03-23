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
	Position  string `json:"-" gorm:"default:'{}'"`
	Icon      string `json:"icon" gorm:"default:'􀛮'"` // The icon to represent this entity
	Frequency int    `json:"-" gorm:"default:3000"`
	Neural    string `json:"-" gorm:"default:'inactive'"` // Parent Module name
	Predicted string `gorm:"-" json:"predicted"`          // scalar
}
