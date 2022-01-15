// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"fmt"
	"sync"
	"udap/internal/bond"
	"udap/internal/models"
	"udap/internal/store"
)

type Generic struct {
	store sync.Map
}

type Modules struct {
	PolyBuffer
	bond *bond.Bond
}

func LoadModules() (m *Modules) {
	m = &Modules{}
	m.data = sync.Map{}
	return m
}

func (m *Modules) Handle(event bond.Msg) (res any, err error) {
	switch o := event.Operation; o {
	case "register":
		return m.register(event)
	default:
		return nil, fmt.Errorf("invalid operation '%s'", o)
	}
}

func (m *Modules) register(event bond.Msg) (res any, err error) {
	module := event.Body.(*models.Module)
	err = module.Emplace()
	if err != nil {
		return nil, err
	}
	m.Set(module.Id, module)
	return nil, nil
}

func (m *Modules) Register(module *models.Module) error {
	err := module.Emplace()
	if err != nil {
		return err
	}
	m.Set(module.Id, module)
	return nil
}

func (m *Modules) FetchAll() {
	var modules []*models.Module
	store.DB.Model(&models.Module{}).Find(&modules)
	for _, module := range modules {
		m.set(module.Id, module)
	}
}

// Pull is the level at which this service needs to run
func (m *Modules) Pull() {
	for _, k := range m.Keys() {
		err := m.get(k)
		if err != nil {
			return
		}
	}
}

func (m *Modules) Compile() (es []models.Module, err error) {
	for _, k := range m.Keys() {
		ea := m.get(k).(*models.Module)
		es = append(es, *ea)
	}
	return es, err
}

func (m *Modules) Find(name string) *models.Module {
	get := m.get(name)
	return get.(*models.Module)
}

func (m *Modules) Set(id string, module *models.Module) {
	m.set(id, module)
}
