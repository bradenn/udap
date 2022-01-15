// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"encoding/json"
	"sync"
	"udap/internal/bond"
	"udap/internal/models"
	"udap/internal/store"
)

type Entities struct {
	PolyBuffer
}

func (e *Entities) Handle(event bond.Msg) (any, error) {
	switch event.Operation {
	case "compile":
		return e.compile(event)
	case "logs":
		return e.logs(event)
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
	case "pull":
		return e.pull(event)
	case "find":
		return e.find(event)
	case "neural":
		return e.neural(event)
	case "predict":
		return e.predict(event)
	case "attribute":
		return e.attribute(event)
	case "state":
		return e.state(event)

	}
	return nil, nil
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
		err = find.Pull()
		if err != nil {
			continue
		}
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

func (e *Entities) Config(id string, data string) (res any, err error) {
	entity := e.Find(id)
	err = entity.ChangeConfig(data)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Entities) pull(event bond.Msg) (res any, err error) {
	res, err = e.Pull(event.Id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *Entities) Pull(id string) (res string, err error) {
	entity := e.Find(id)
	err = entity.Pull()
	if err != nil {
		return "", err
	}
	return "", nil
}

func (e *Entities) State(id string, state string) (res any, err error) {
	entity := e.Find(id)
	err = entity.Push(state)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Entities) state(event bond.Msg) (res any, err error) {
	_, err = e.State(event.Id, event.Payload)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func LoadEntities() (m *Entities) {
	m = &Entities{}
	m.raw = map[string]any{}
	m.data = sync.Map{}
	m.FetchAll()
	return m
}

func (e *Entities) FetchAll() {
	var entities []*models.Entity
	store.DB.Model(&models.Entity{}).Find(&entities)
	for _, entity := range entities {
		err := entity.Pull()
		if err != nil {
		}
		e.set(entity.Id, entity)
	}
}

func (e *Entities) PullAll() {
	for _, k := range e.Keys() {
		_, err := e.Pull(k)
		if err != nil {
			return
		}
	}
}

func (e *Entities) Find(name string) *models.Entity {
	return e.get(name).(*models.Entity)
}

func (e *Entities) Set(id string, entity *models.Entity) {
	e.set(id, entity)
}
func (e *Entities) Logs(id string) ([]models.Log, error) {
	entity := e.Find(id)
	return entity.Logs()
}
func (e *Entities) logs(msg bond.Msg) (any, error) {
	var lgs []models.Log

	err := store.DB.Model(&models.Log{}).Where("entity_id = ? AND cct <> 0", msg.Id).Find(&lgs).Error
	return lgs, err
}

func (e *Entities) attribute(event bond.Msg) (any, error) {
	entity := e.Find(event.Id)
	err := entity.Push(event.Payload)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
