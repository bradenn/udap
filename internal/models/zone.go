// Copyright (c) 2021 Braden Nicholson

package models

import (
	"time"
	"udap/internal/log"
	"udap/internal/store"
)

type Zone struct {
	store.Persistent
	Name     string `json:"name"`
	Entities string `json:"entities"`
}

func (z *Zone) Emplace() (err error) {
	z.UpdatedAt = time.Now()
	err = store.DB.Model(&Zone{}).Where("id = ?", z.Id).FirstOrCreate(z).Error
	if err != nil {
		return err
	}
	return nil
}

func (z *Zone) FetchAll() []Zone {
	var zones []Zone
	log.Log("Fetching")
	err := store.DB.Table("zones").Find(&zones).Error
	if err != nil {
		return nil
	}
	return zones
}

func (z *Zone) Update() error {
	err := store.DB.Where("id = ?", z.Id).Save(&z).Error
	return err
}

func NewZone() Zone {
	zone := Zone{}
	return zone
}
