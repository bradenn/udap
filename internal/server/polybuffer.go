// Copyright (c) 2021 Braden Nicholson

package server

import (
	"sync"
	"udap/internal/cache"
	"udap/internal/models"
	"udap/internal/store"
	"udap/pkg/plugin"
)

type PolyBuffer struct {
	data sync.Map
}

func LoadEntities() (m *Entities) {
	m = &Entities{}
	m.data = sync.Map{}
	m.FetchAll()
	return m
}

func LoadEndpoints() (m *Endpoints) {
	m = &Endpoints{}
	m.data = sync.Map{}
	m.FetchAll()
	return m
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

type Entity struct {
	*models.Entity
}

func (e *Entities) FetchAll() {
	var entities []Entity
	store.DB.Model(&Entity{}).Find(&entities)
	for _, entity := range entities {
		ln, err := cache.GetLn(entity.Path(), "state")
		if err != nil {
			mono := models.Mono{Value: 0}
			entity.State = mono.Marshal()
			_ = cache.PutLn(entity.State, entity.Path(), "state")
			e.Set(entity.Path(), entity)
		} else {
			entity.State = []byte(ln.(string))
			e.Set(entity.Path(), entity)
		}

	}
}

func (e *Entities) Compile() (es []Entity, err error) {
	for _, k := range e.Keys() {
		ea := e.get(k).(Entity)
		if ea.Connected() {
			err = ea.Poll()
			if err != nil {
				return nil, err
			}
		}
		es = append(es, ea)
	}
	return es, err
}

func (e *Entities) Find(name string) Entity {
	return e.get(name).(Entity)
}

func (e *Entities) register(name string, entity Entity) {
	e.set(name, entity)
}

func (e *Entities) Set(name string, entity Entity) {
	e.set(name, entity)
}

type Endpoints struct {
	PolyBuffer
}

func (e *Endpoints) FetchAll() {
	e.data = sync.Map{}
	var endpoints []models.Endpoint
	store.DB.Find(&endpoints)
	for _, endpoint := range endpoints {
		e.Set(endpoint.Id, endpoint)
	}
}

func (e *Endpoints) Compile() (endpoints []models.Endpoint, err error) {
	for _, s := range e.Keys() {
		endpoint := e.get(s).(models.Endpoint)
		endpoints = append(endpoints, endpoint)
	}
	return endpoints, err
}

func (e *Endpoints) Find(name string) models.Endpoint {
	dat := e.get(name)
	if dat == nil {
		endpoint := models.Endpoint{}
		err := endpoint.Fetch()
		if err != nil {
			return models.Endpoint{}
		}
		return endpoint
	}
	return dat.(models.Endpoint)
}

func (e *Endpoints) Set(name string, endpoint models.Endpoint) {
	e.set(name, endpoint)
}

func (e *Endpoints) Save(endpoint *models.Endpoint) {
	store.DB.Model(&models.Endpoint{}).Save(endpoint)
}
