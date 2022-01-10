// Copyright (c) 2021 Braden Nicholson

package controller

import (
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

func (d *Devices) register(event bond.Msg) (res any, err error) {
	device := event.Body.(*models.Device)
	err = device.Emplace()
	if err != nil {
		return nil, err
	}
	d.Set(device.Id, device)
	return nil, nil
}

func LoadDevices() (m *Devices) {
	m = &Devices{}
	m.raw = map[string]any{}
	m.FetchAll()
	return m
}

func (d *Devices) FetchAll() {
	var devices []*models.Device
	store.DB.Model(&models.Device{}).Find(&devices)
	for _, device := range devices {
		d.set(device.Id, device)
	}
}

// Pull is the level at which this service needs to run
func (d *Devices) Pull() {
	for _, k := range d.Keys() {
		err := d.get(k)
		if err != nil {
			return
		}
	}
}

func (d *Devices) Compile() (es []models.Device, err error) {
	for _, k := range d.Keys() {
		ea := d.get(k).(*models.Device)
		es = append(es, *ea)
	}
	return es, err
}

func (d *Devices) compile(msg bond.Msg) (res any, err error) {
	var devices []models.Device
	for _, s := range d.Keys() {
		device := d.Find(s)
		devices = append(devices, *device)
	}
	return devices, nil
}

func (d *Devices) Find(name string) *models.Device {
	return d.get(name).(*models.Device)
}

func (d *Devices) Set(id string, entity *models.Device) {
	d.set(id, entity)
}
