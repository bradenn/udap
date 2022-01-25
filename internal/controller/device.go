// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"sync"
	"udap/internal/bond"
	"udap/internal/log"
	"udap/internal/models"
	"udap/internal/store"
)

type Devices struct {
	PolyBuffer
	Observable
}

func (d *Devices) Handle(event bond.Msg) (res any, err error) {
	switch event.Operation {
	}
	return nil, nil
}

func (d *Devices) Compile() (res []models.Device, err error) {
	for _, s := range d.Keys() {
		device := d.Find(s)
		if device == nil {
			continue
		}

		res = append(res, *device)
	}
	return res, nil
}

func (d *Devices) FetchAll() {
	var devices []models.Device
	err := store.DB.Model(&models.Device{}).Find(&devices).Error
	if err != nil {
		log.Err(err)
		return
	}
	for _, device := range devices {

		d.set(device.Id, &device)
		d.emit(device.Id, &device)
	}
}

func (d *Devices) Find(name string) *models.Device {
	return d.get(name).(*models.Device)
}

func LoadDevices() (m *Devices) {
	m = &Devices{}
	m.data = sync.Map{}
	m.Run()
	m.FetchAll()
	return m
}

func (d *Devices) Register(device models.Device) (res *models.Device, err error) {
	err = device.Emplace()
	if err != nil {
		return nil, err
	}
	d.set(device.Id, &device)
	d.emit(device.Id, &device)
	return nil, nil
}

func (d *Devices) Set(id string, device *models.Device) {
	d.set(id, device)
	d.emit(device.Id, &device)
}
