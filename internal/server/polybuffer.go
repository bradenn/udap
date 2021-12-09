// Copyright (c) 2021 Braden Nicholson

package server

import (
	"sync"
	"udap/internal/models"
	"udap/internal/store"
	"udap/pkg/plugin"
)

type PolyBuffer struct {
	data sync.Map
}

func (p *PolyBuffer) set(name string, data interface{}) {
	p.data.Store(name, data)
}

func (p *PolyBuffer) get(name string) interface{} {
	load, ok := p.data.Load(name)
	if ok {
		return load
	}
	return struct{}{}
}

func (p *PolyBuffer) Keys() []string {
	var s []string
	p.data.Range(func(key, value interface{}) bool {
		s = append(s, key.(string))
		return true
	})
	return s
}

type Modules struct {
	PolyBuffer
}

func (m *Modules) Find(name string) plugin.UdapPlugin {
	return m.get(name).(plugin.UdapPlugin)
}

func (m *Modules) Set(name string, module plugin.UdapPlugin) {
	m.set(name, module)
}

type Entities struct {
	PolyBuffer
}

func (e *Entities) FetchAll() {
	var entities []models.Entity
	store.DB.Model(&models.Entity{}).Find(&entities)
	for _, entity := range entities {
		e.Set(entity.Id, &entity)
	}
}

func (e *Entities) Find(name string) *models.Entity {
	return e.get(name).(*models.Entity)
}

func (e *Entities) register(name string, entity *models.Entity) {
	e.set(name, entity)
}

func (e *Entities) Set(name string, entity *models.Entity) {
	e.set(name, entity)
}

type Endpoints struct {
	PolyBuffer
}

func (e *Endpoints) FetchAll() {
	var endpoints []models.Endpoint
	store.DB.Model(&models.Endpoint{}).Find(&endpoints)
	for _, endpoint := range endpoints {
		e.Set(endpoint.Id, &endpoint)
	}
}

func (e *Endpoints) Find(name string) *models.Endpoint {
	dat := e.get(name)
	if dat == nil {
		endpoint := &models.Endpoint{}
		err := endpoint.Fetch()
		if err != nil {
			return nil
		}
		return endpoint
	}
	return dat.(*models.Endpoint)
}

func (e *Endpoints) Set(name string, endpoint *models.Endpoint) {
	e.set(name, endpoint)
}

func (e *Endpoints) Save(endpoint *models.Endpoint) {
	store.DB.Model(&models.Endpoint{}).Save(endpoint)
}
