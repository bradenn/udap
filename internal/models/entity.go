// Copyright (c) 2021 Braden Nicholson

package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"udap/internal/log"
	"udap/internal/store"
)

type State json.RawMessage

type Entity struct {
	store.Persistent
	Name   string `json:"name"`
	Alias  string `json:"alias"`
	Type   string `json:"type"`
	Module string `json:"module"`
	Locked bool   `json:"locked"`
	State  State  `gorm:"-" json:"state"`
	Live   bool   `gorm:"-" json:"live"`
	tx     Tx
	rx     Rx
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
	e.Live = true
	err := e.get()
	if err != nil {
		return err
	}
	e.tx = tx

	return nil
}

func (e *Entity) Connected() bool {
	return e.rx != nil && e.tx != nil
}

func (e *Entity) send(payload State) error {
	if e.tx == nil {
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
	if e.rx == nil {
		log.Err(fmt.Errorf("yah brah, no rx kyyyyd"))
		mono := Mono{
			Value: 0.0,
		}
		e.State = mono.Marshal()
		return nil
	}
	e.State = e.rx()
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

func (e *Entity) GetState() (State, error) {
	err := e.get()
	if err != nil {
		return nil, err
	}
	return e.State, nil
}

func (e *Entity) SetState(state State) error {
	err := e.send(state)
	if err != nil {
		return err
	}
	e.State = state
	return nil
}
