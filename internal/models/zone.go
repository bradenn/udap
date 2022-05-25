// Copyright (c) 2021 Braden Nicholson

package models

import (
	"time"
	"udap/internal/store"
)

type Zone struct {
	store.Persistent
	Name     string   `json:"name"`
	Entities []Entity `json:"entities" gorm:"many2many:zone_entities;"`
	User     string   `json:"user"`
}

// Emplace will Find or Create a zone based on its id.
func (z *Zone) Emplace() (err error) {
	z.UpdatedAt = time.Now()
	err = store.DB.Model(&Zone{}).FirstOrCreate(z).Error
	if err != nil {
		return err
	}
	return nil
}

func (z *Zone) FetchAll() (err error, zones []Zone) {
	if err = store.DB.Table("zones").Preload("Entities").Find(&zones).Error; err != nil {
		return err, nil
	}
	return nil, zones
}

func (z *Zone) Update() error {
	return store.DB.Where("id = ?", z.Id).Save(&z).Error
}

func (z *Zone) Restore() error {
	z.Deleted = false
	return z.Update()
}

// Delete marks the zone as deleted and discontinues its function
func (z *Zone) Delete() error {
	z.Deleted = true
	return z.Update()
}
