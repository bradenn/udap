// Copyright (c) 2022 Braden Nicholson

package zone

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
)

type zoneRepo struct {
	generic.Store[domain.Zone]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.ZoneRepository {
	return &zoneRepo{
		db:    db,
		Store: generic.NewStore[domain.Zone](db),
	}
}
