// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"encoding/json"
	"sync"
	"udap/internal/bond"
	"udap/internal/log"
	"udap/internal/models"
)

type Entities struct {
	PolyBuffer
	Observable
}

func (e *Entities) Handle(event bond.Msg) (any, error) {
	switch event.Operation {
	case "register":
		return e.register(event)
	case "rename": // Alias
		return e.rename(event)
	case "lock":
		return e.lock(event)
	case "unlock":
		return e.unlock(event)
	case "icon":
		return e.icon(event)
	case "neural":
		return e.neural(event)
	case "predict":
		return e.predict(event)
	}
	return nil, nil
}

func (e *Entities) neural(event bond.Msg) (res any, err error) {
	entity := e.Find(event.Id)
	ref := e.Parse(event.Payload)
	err = entity.ChangeNeural(ref.Neural)
	if err != nil {
		return nil, err
	}
	return nil, err
}

func (e *Entities) EmitAll() (err error) {

	for _, k := range e.Keys() {
		find := e.Find(k)
		e.emit(k, find)
	}

	return nil
}

func (e *Entities) register(event bond.Msg) (any, error) {
	entity := e.Cast(event.Body)

	_, err := e.Register(entity)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Entities) find(event bond.Msg) (res any, err error) {
	entity := e.Find(event.Id)
	return entity, nil
}

func (e *Entities) Suggest(id string, body string) (res any, err error) {
	entity := e.Find(id)
	err = entity.Suggest(body)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Entities) predict(event bond.Msg) (res any, err error) {
	entity := e.Find(event.Id)
	err = entity.Suggest(string(event.Body.(json.RawMessage)))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Entities) Cast(body any) *models.Entity {
	return body.(*models.Entity)
}

func (e *Entities) Parse(body string) models.Entity {
	entity := models.Entity{}
	err := json.Unmarshal([]byte(body), &entity)
	if err != nil {
		return models.Entity{}
	}
	return entity
}

func (e *Entities) Register(entity *models.Entity) (res *models.Entity, err error) {
	log.Event("Entity '%s' registered.", entity.Name)
	err = entity.Emplace()
	if err != nil {
		return nil, err
	}
	e.Set(entity.Id, entity)
	return entity, nil
}

func (e *Entities) rename(event bond.Msg) (res any, err error) {
	ref := e.Cast(event.Body)
	_, err = e.Rename(event.Id, ref.Alias)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Entities) Rename(id string, name string) (res any, err error) {
	entity := e.Find(id)
	err = entity.Rename(name)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Entities) lock(event bond.Msg) (res any, err error) {
	entity := e.Find(event.Id)
	err = entity.Lock()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Entities) unlock(event bond.Msg) (res any, err error) {
	entity := e.Find(event.Id)
	err = entity.Unlock()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Entities) icon(event bond.Msg) (res any, err error) {
	entity := e.Find(event.Id)
	ee := e.Parse(event.Payload)
	err = entity.ChangeIcon(ee.Icon)
	if err != nil {
		return nil, err
	}
	e.Set(event.Id, entity)
	return nil, nil
}

func (e *Entities) Config(id string, data string) (res any, err error) {
	entity := e.Find(id)
	err = entity.ChangeConfig(data)
	if err != nil {
		return nil, err
	}
	e.Set(id, entity)
	return nil, nil
}

func LoadEntities() (m *Entities) {
	m = &Entities{}
	m.raw = map[string]any{}
	m.data = sync.Map{}
	m.Run()
	return m
}

func (e *Entities) Find(name string) *models.Entity {
	en := e.get(name).(*models.Entity)
	return en
}

func (e *Entities) Set(id string, entity *models.Entity) {
	e.set(id, entity)
	e.emit(id, entity)
}
