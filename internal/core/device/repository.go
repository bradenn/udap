// Copyright (c) 2022 Braden Nicholson

package device

import (
	"gorm.io/gorm"
	"udap/internal/core/generic"
)

type deviceRepo struct {
	generic.Store[Device]
	db *gorm.DB
}

func newRepository(db *gorm.DB) Repository {
	return &deviceRepo{
		db:    db,
		Store: generic.NewStore[Device](db),
	}
}

func (m *deviceRepo) FindOrCreate(device *Device) error {
	return m.db.FirstOrCreate(device, "mac = ?", device.Mac).Error
}
