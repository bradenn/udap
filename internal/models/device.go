// Copyright (c) 2021 Braden Nicholson

package models

import (
	"time"
	"udap/internal/log"
	"udap/internal/store"
)

type Device struct {
	store.Persistent
	NetworkId string `json:"networkId" gorm:"-"`
	EntityId  string `json:"entityId" gorm:"-"`
	Name      string `json:"name"`
	Hostname  string `json:"hostname"`
	Mac       string `json:"mac"`
	Ipv4      string `json:"ipv4"`
	Ipv6      string `json:"ipv6"`
}

// Emplace gets a module from its path
func (d *Device) Emplace() (err error) {
	d.UpdatedAt = time.Now()
	err = store.DB.Model(&Device{}).Where("mac = ?", d.Mac).FirstOrCreate(&d).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) FetchAll() []Device {
	var devices []Device
	log.Log("Fetching")
	err := store.DB.Table("devices").Find(&devices).Error
	if err != nil {
		return nil
	}
	return devices
}

func (d *Device) Update() error {
	err := store.DB.Where("mac = ?", d.Mac).Save(&d).Error
	return err
}

func NewDevice() Device {
	device := Device{}
	return device
}
