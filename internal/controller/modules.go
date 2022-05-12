// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"fmt"
	"sync"
	"udap/internal/bond"
	"udap/internal/models"
	"udap/internal/store"
)

type Modules struct {
	PolyBuffer
	Observable
}

func LoadModules() (m *Modules) {
	m = &Modules{}
	m.data = sync.Map{}
	m.raw = map[string]interface{}{}
	m.Run()
	m.FetchAll()
	return m
}

func (m *Modules) EmitAll() (err error) {

	for _, k := range m.Keys() {
		find := m.Find(k)
		m.emit(k, find)
	}

	return nil
}

func (m *Modules) Handle(event bond.Msg) (res interface{}, err error) {
	switch o := event.Operation; o {
	case "register":
		return m.register(event)
	case "enabled":
		return m.enable(event)
	default:
		return nil, fmt.Errorf("invalid operation '%s'", o)
	}
}

func (m *Modules) enable(event bond.Msg) (res interface{}, err error) {
	id := event.Id
	err = m.Enabled(id, event.Payload == "true")
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (m *Modules) register(event bond.Msg) (res interface{}, err error) {
	module := event.Body.(*models.Module)
	err = module.Emplace()
	if err != nil {
		return nil, err
	}
	m.Set(module.Id, module)
	return nil, nil
}

func (m *Modules) Register(module models.Module) (string, error) {

	err := module.Emplace()
	if err != nil {
		return "", err
	}

	m.Set(module.Id, &module)
	return module.Id, nil
}

func (m *Modules) State(id string, state string) error {
	find := m.Find(id)
	find.State = state

	err := find.Update()
	if err != nil {
		return err
	}

	m.Set(find.Id, find)
	return nil
}

func (m *Modules) Enabled(id string, enabled bool) error {
	find := m.Find(id)
	find.Enabled = enabled

	err := find.Update()
	if err != nil {
		return err
	}

	m.Set(find.Id, find)
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
	m.emit(id, module)
}
