// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type Zone struct {
	common.Persistent
	Name     string   `json:"name"`
	Entities []Entity `json:"entities" gorm:"many2many:zone_entities;"`
	Pinned   bool     `json:"pinned"`
	User     string   `json:"user"`
}
