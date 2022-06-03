// Copyright (c) 2022 Braden Nicholson

package device

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
)

type deviceRepo struct {
	generic.Store[domain.Device]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.DeviceRepository {
	return &deviceRepo{
		db:    db,
		Store: generic.NewStore[domain.Device](db),
	}
}
