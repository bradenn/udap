// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"udap/internal/bond"
	"udap/internal/models"
	"udap/internal/store"
)

type Networks struct {
	PolyBuffer
}

func (d *Networks) Handle(event bond.Msg) (res any, err error) {
	switch event.Operation {
	case "register":
		return d.register(event)
	case "compile":
		return d.compile(event)
	}
	return nil, nil
}

func (d *Networks) compile(msg bond.Msg) (res any, err error) {
	var networks []models.Network
	for _, s := range d.Keys() {
		network := d.Find(s)
		networks = append(networks, *network)
	}
	return networks, nil
}

func (d *Networks) register(event bond.Msg) (res any, err error) {
	network := event.Body.(*models.Network)
	err = network.Emplace()
	if err != nil {
		return nil, err
	}
	d.Set(network.Id, network)
	return nil, nil
}

func LoadNetworks() (m *Networks) {
	m = &Networks{}
	m.raw = map[string]any{}
	m.FetchAll()
	return m
}

func (d *Networks) FetchAll() {
	var networks []*models.Network
	store.DB.Model(&models.Network{}).Find(&networks)
	for _, network := range networks {
		d.set(network.Id, network)
	}
}

// Pull is the level at which this service needs to run
func (d *Networks) Pull() {
	for _, k := range d.Keys() {
		err := d.get(k)
		if err != nil {
			return
		}
	}
}

func (d *Networks) Compile() (es []models.Network, err error) {
	for _, k := range d.Keys() {
		ea := d.get(k).(*models.Network)
		es = append(es, *ea)
	}
	return es, err
}

func (d *Networks) Find(name string) *models.Network {
	return d.get(name).(*models.Network)
}

func (d *Networks) Set(id string, entity *models.Network) {
	d.set(id, entity)
}
