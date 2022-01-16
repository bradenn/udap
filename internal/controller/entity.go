// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"encoding/json"
	"fmt"
	"sync"
	"udap/internal/bond"
	"udap/internal/cache"
	"udap/internal/models"
)

type Entities struct {
	PolyBuffer
}

func (e *Entities) Handle(event bond.Msg) (any, error) {
	switch event.Operation {
	case "compile":
		return e.compile(event)
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
	case "Find":
		return e.find(event)
	case "neural":
		return e.neural(event)
	case "predict":
		return e.predict(event)
	case "attribute":
		return e.attribute(event)
	}
	return nil, nil
}

func (e *Entities) Observe(id string, fn func(attribute models.Entity) error) {
	cache.WatchFn(fmt.Sprintf("entity.%s", id), func(s string) error {
		entity := models.Entity{}
		err := json.Unmarshal([]byte(s), &entity)
		if err != nil {
			return err
		}
		err = fn(entity)
		if err != nil {
			return err
		}
		return nil
	})
}

func (e *Entities) compile(event bond.Msg) (res any, err error) {
	return e.Compile()
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

func (e *Entities) Compile() (res []models.Entity, err error) {

	var entities []models.Entity
	for _, k := range e.Keys() {
		find := e.get(k).(*models.Entity)
		entities = append(entities, *find)
	}

	return entities, nil
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
	return nil, nil
}

func (e *Entities) Attribute(id string, key string, value string) (res any, err error) {
	err = cache.PutLn(value, "entity", id, "attribute", key)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Entities) attribute(msg bond.Msg) (res any, err error) {
	a := models.Attribute{}
	err = json.Unmarshal([]byte(msg.Payload), &a)
	if err != nil {
		return nil, err
	}
	_, err = e.Attribute(a.Entity, a.Key, a.Value)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (e *Entities) Config(id string, data string) (res any, err error) {
	entity := e.Find(id)
	err = entity.ChangeConfig(data)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func LoadEntities() (m *Entities) {
	m = &Entities{}
	m.raw = map[string]any{}
	m.data = sync.Map{}
	return m
}

func (e *Entities) Find(name string) *models.Entity {
	en := e.get(name).(*models.Entity)
	return en
}

func (e *Entities) Set(id string, entity *models.Entity) {
	e.set(id, entity)
}
