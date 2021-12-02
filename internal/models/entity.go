// Copyright (c) 2021 Braden Nicholson

package models

import (
	"udap/internal/log"
	"udap/internal/store"
)

const (
	SWITCH = "switch"
	DIMMER = "DIMMER"
	RGB    = "rgb"
	RGBW   = "rgbw"
	RGBCCT = "rgbcct"
)

type Entity struct {
	store.Persistent
	Name     string `json:"name"`
	Type     string `json:"type"`
	ModuleId string `json:"moduleId"`
	State    string `gorm:"-" json:"state"`
}

type Dimmer struct {
	Light
}

type Light struct {
	Entity
}

func (l *Light) TurnOn() error {
	return nil
}

func (l *Light) TurnOff() error {
	return nil
}

func GetEntities() ([]Entity, error) {
	var entities []Entity
	err := store.DB.Model(&Entity{}).Find(&entities).Error
	if err != nil {
		log.Err(err)
	}
	for i, entity := range entities {
		state, err := store.GetLn("entity", entity.Id, "state")
		if err != nil {
			err = store.PutLn("off", "entity", entity.Id, "state")
		}
		entities[i].State = state
	}
	return entities, nil
}

// Fetch gets a module from its path
func (e *Entity) Fetch() (err error) {
	err = store.DB.Model(&Entity{}).Where("id = ?", e.Id).FirstOrCreate(e).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) SetState(state string) {
	err := store.PutLn(state, "entity", e.Id, "state")
	if err != nil {
		return
	}
}
