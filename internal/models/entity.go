// Copyright (c) 2021 Braden Nicholson

package models

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"udap/internal/cache"
	"udap/internal/store"
)

type IController interface {
	Resolve() (string, error)
	Request(state string) error
}

type Controller struct {
}

func (c Controller) Resolve() (string, error) {
	return "", nil
}

func (c Controller) Request(state string) error {
	return nil
}

type Device interface {
}

type Entity struct {
	store.Persistent
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Module     string      `json:"module"`
	State      interface{} `gorm:"-" json:"state"`
	fnSetState func(state interface{}) error
	fnGetState func() interface{}
}

func NewSwitch(name string, module string) *Entity {
	e := Entity{
		Name:   name,
		Type:   "switch",
		Module: module,
		State:  "off",
	}
	return &e
}

func NewDimmer(name string, module string) *Entity {
	e := Entity{
		Name:   name,
		Type:   "dimmer",
		Module: module,
		State:  "0",
	}
	return &e
}

// Find attempts to locate
func (e *Entity) Find() error {
	err := store.DB.Where("name = ? AND module = ?", e.Name, e.Module).First(&e).Error
	return err
}

// Path attempts to locate
func (e *Entity) Path() string {

	return strings.ToLower(fmt.Sprintf("%s.%s", e.Module, e.Name))
}

// Find attempts to locate
func (e *Entity) SetStateHandler(fn func(state interface{}) error) {
	e.fnSetState = fn
}

// Find attempts to locate
func (e *Entity) SetStateReceiver(fn func() interface{}) {
	e.fnGetState = fn
}

// Deserialize attempts to locate
func (e *Entity) Deserialize(data map[string]interface{}) error {
	e.Id = data["id"].(string)
	e.Name = data["name"].(string)
	e.Type = data["type"].(string)
	e.Module = data["module"].(string)
	e.State = data["state"].(string)
	return nil
}

// Serialize attempts to locate
func (e *Entity) Serialize() map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = e.Id
	data["name"] = e.Name
	data["type"] = e.Type
	data["module"] = e.Module
	data["state"] = e.State
	return data
}

func (e *Entity) Insert() error {

	return nil
}

func (e *Entity) Delete() error {
	err := store.DB.Where("name = ? AND module = ?", e.Name, e.Module).Delete(&e).Error
	return err
}

func (e *Entity) Update() error {
	err := store.DB.Where("id = ?", e.Id).Save(&e).Error
	return err
}

// Emplace gets a module from its path
func (e *Entity) Emplace() (err error) {
	err = store.DB.Model(&Entity{}).Where("name = ? AND module = ?", e.Name, e.Module).FirstOrCreate(&e).Error
	if err != nil {
		return err
	}
	return nil
}

// Fetch gets a module from its path
func (e *Entity) Fetch() (err error) {
	err = store.DB.Model(&Entity{}).Where("name = ? AND module = ?", e.Name, e.Module).First(&e).Error
	if err != nil {
		return err
	}
	return nil
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

	}
	return entities, nil
}

func (e *Entity) updateState() error {
	ln, err := cache.GetLn(e.Module, e.Name, "state")
	e.State = ln
	if err != nil {
		switch e.Type {
		case "switch":
			e.State = false
			_ = cache.PutLn(false, e.Module, e.Name, "state")
		case "dimmer":
			e.State = 0
			_ = cache.PutLn(0, e.Module, e.Name, "state")
		}
	}
	if e.fnSetState != nil {
		err = e.fnSetState(e.State)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Entity) Poll() error {

	return e.updateState()
}

func (e *Entity) AfterFind(_ *gorm.DB) error {
	return e.updateState()
}

func (e *Entity) SetState() error {
	if e.fnSetState == nil {
		return fmt.Errorf("fnSetState not set")
	}

	err := cache.PutLn(e.State, e.Module, e.Name, "state")
	if err != nil {

	}

	err = e.fnSetState(e.State)
	if err != nil {
		return err
	}

	return nil
}
