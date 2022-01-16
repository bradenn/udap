// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"sync"
	"udap/internal/bond"
	"udap/internal/models"
	"udap/internal/store"
)

type Devices struct {
	PolyBuffer
}

func (d *Devices) Handle(event bond.Msg) (res any, err error) {
	switch event.Operation {
	case "register":
		return d.register(event)
	case "compile":
		return d.compile(event)
	}
	return nil, nil
}

func (d *Devices) Compile() (es []models.Device, err error) {
	var devices []models.Device
	for _, s := range d.Keys() {
		device := d.Find(s)
		devices = append(devices, *device)
	}
	return es, err
}

func (d *Devices) FetchAll() {
	var devices []*models.Device
	store.DB.Model(&models.Device{}).Find(&devices)
	for _, device := range devices {
		d.set(device.Id, device)
	}
}

func (d *Devices) Find(name string) *models.Device {
	return d.get(name).(*models.Device)
}

func LoadDevices() (m *Devices) {
	m = &Devices{}
	m.data = sync.Map{}
	m.FetchAll()
	return m
}

func (d *Devices) Pull() {
	for _, k := range d.Keys() {
		err := d.get(k)
		if err != nil {
			return
		}
	}
}

func (d *Devices) Register(device *models.Device) (res *models.Device, err error) {
	err = device.Emplace()
	if err != nil {
		return nil, err
	}
	d.Set(device.Id, device)
	return nil, nil
}

func (d *Devices) Set(id string, entity *models.Device) {
	d.set(id, entity)
}

// Bond

func (d *Devices) compile(msg bond.Msg) (res any, err error) {
	return d.Compile()
}

func (d *Devices) register(event bond.Msg) (res any, err error) {
	device := event.Body.(*models.Device)
	return d.Register(device)
}
