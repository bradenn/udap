// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"udap/internal/bond"
	"udap/internal/models"
	"udap/internal/store"
)

type Modules struct {
	PolyBuffer
	bond *bond.Bond
}

func (m *Modules) Handle(event bond.Msg) (res any, err error) {
	switch event.Operation {
	case "register":
		return m.register(event)

	}
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

func (m *Modules) register(event bond.Msg) (res any, err error) {
	module := event.Body.(*models.Module)
	err = module.Emplace()
	if err != nil {
		return nil, err
	}
	m.Set(module.Id, module)
	return nil, nil
}

func LoadModules() (m *Modules) {
	m = &Modules{}
	m.raw = map[string]any{}
	return m
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
	return m.get(name).(*models.Module)
}

func (m *Modules) Set(id string, module *models.Module) {
	m.set(id, module)
}
