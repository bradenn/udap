// Copyright (c) 2022 Braden Nicholson

package models

import (
	"fmt"
	"strconv"
	"udap/internal/store"
)

type FuncPut func(Attribute) error
type FuncGet func(Attribute) (string, error)

type Attribute struct {
	store.Persistent
	EntityId string `json:"entity_id"`
	Key      string `json:"key"`
	Value    string `json:"value" gorm:"-"`
	put      FuncPut
	get      FuncGet
}

func (a *Attribute) Get() error {
	if a.get == nil {
		return fmt.Errorf("attribute disconnected from entity module")
	}
	res, err := a.get(*a)
	if err != nil {
		return err
	}
	a.Value = res
	return nil
}

func (a *Attribute) Put(body string) error {
	a.Value = body
	err := a.put(*a)
	if err != nil {
		return err
	}
	a.Value = body
	return nil
}

func (a *Attribute) FnPut(put FuncPut) {
	a.put = put
}

func (a *Attribute) FnGet(get FuncGet) {
	a.get = get
}

func (a *Attribute) AsInt() int {
	parsed, err := strconv.Atoi(a.Value)
	if err != nil {
		return 0
	}
	return parsed
}

func (a *Attribute) AsFloat() float64 {
	parsed, err := strconv.ParseFloat(a.Value, 64)
	if err != nil {
		return 0.0
	}
	return parsed
}

func (a *Attribute) AsBool() bool {
	parsed, err := strconv.ParseBool(a.Value)
	if err != nil {
		return false
	}
	return parsed
}
