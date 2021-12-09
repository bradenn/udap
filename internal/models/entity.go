// Copyright (c) 2021 Braden Nicholson

package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"udap/internal/store"
)

type State json.RawMessage

type Entity struct {
	store.Persistent
	Name   string `json:"name"`
	Type   string `json:"type"`
	Module string `json:"module"`
	Locked bool   `json:"locked"`
	State  State  `gorm:"-" json:"state"`
	live   bool
	tx     Tx
	rx     Rx
}

// Find attempts to locate
func (e *Entity) Live() bool {
	return e.live
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

type Tx func(state State) error
type Rx func() State

// Handlers attempts to locate
func (e *Entity) Handlers(tx Tx, rx Rx) error {
	e.rx = rx
	err := e.get()
	if err != nil {
		return err
	}
	e.tx = tx
	e.live = true
	return nil
}

func (e *Entity) send(payload State) error {
	if e.tx == nil || !e.live {
		return fmt.Errorf("entity '%s' is not connected to its parent module '%s'", e.Name, e.Module)
	}
	if e.Locked {
		return fmt.Errorf("entity '%s' is locked. Unlock it before making changes", e.Name)
	}
	err := e.tx(payload)
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) get() error {
	data := e.rx()
	e.State = data
	return nil
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

func (e *Entity) Poll() error {
	return e.get()
}

func (e *Entity) AfterFind(_ *gorm.DB) error {
	return nil
}

func (e *Entity) Push(state State) error {
	err := e.send(state)
	if err != nil {
		return err
	}
	e.State = state
	return nil
}
