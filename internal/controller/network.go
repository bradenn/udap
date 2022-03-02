// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"sync"
	"udap/internal/bond"
	"udap/internal/models"
	"udap/internal/store"
)

type Networks struct {
	PolyBuffer
	Observable
}

func (d *Networks) Handle(event bond.Msg) (res any, err error) {
	switch event.Operation {
	case "register":
		return d.register(event)
	}
	return nil, nil
}

func (d *Networks) Compile() (res []models.Network, err error) {
	for _, s := range d.Keys() {
		network := d.Find(s)
		res = append(res, *network)
	}
	return res, nil
}

func (d *Networks) Register(network *models.Network) (res *models.Network, err error) {
	err = network.Emplace()
	if err != nil {
		return nil, err
	}
	d.Set(network.Id, network)
	return nil, nil
}

func (d *Networks) register(event bond.Msg) (res *models.Network, err error) {
	network := event.Body.(*models.Network)

	return d.Register(network)
}

func LoadNetworks() (m *Networks) {
	m = &Networks{}
	m.data = sync.Map{}
	m.raw = map[string]any{}
	m.Run()
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

func (d *Networks) Find(name string) *models.Network {
	return d.get(name).(*models.Network)
}

func (d *Networks) Set(id string, entity *models.Network) {
	d.set(id, entity)
	d.emit(id, entity)
}
