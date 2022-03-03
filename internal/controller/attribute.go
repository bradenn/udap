// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"udap/internal/bond"
	"udap/internal/models"
)

type Attributes struct {
	PolyBuffer
	Observable
}

func (a *Attributes) Handle(event bond.Msg) (res interface{}, err error) {
	switch event.Operation {
	case "request":
		return a.request(event)
	default:
		return nil, fmt.Errorf("operaiton not found")
	}
}

func (a *Attributes) Register(attribute *models.Attribute) error {
	attribute.Id = attribute.Path()
	a.Store(attribute)
	return nil
}

func (a *Attributes) Request(entity string, key string, value string) error {
	attr := models.Attribute{}
	attr.Entity = entity
	attr.Key = key
	attribute := a.Find(attr.Path())
	if attribute == nil {
		return fmt.Errorf("attribute '%s' not found", attr.Id)
	}

	err := attribute.SendRequest(value)
	if err != nil {
		return err
	}

	a.Store(attribute)

	return nil
}

func (a *Attributes) request(event bond.Msg) (interface{}, error) {
	attr := models.Attribute{}
	err := json.Unmarshal([]byte(event.Payload), &attr)
	if err != nil {
		return nil, err
	}
	attribute := a.Find(attr.Path())
	if attribute == nil {
		return nil, fmt.Errorf("attribute '%s' not found", attr.Id)
	}

	err = attribute.SendRequest(attr.Request)
	if err != nil {
		return nil, err
	}

	a.Store(attribute)

	return nil, nil
}

func (a *Attributes) EmitAll() (err error) {

	for _, k := range a.Keys() {
		find := a.Find(k)
		a.emit(k, find)
	}

	return nil
}

func (a *Attributes) Set(entity string, key string, value string) error {
	attr := models.Attribute{}
	attr.Entity = entity
	attr.Key = key

	attribute := a.Find(attr.Path())
	if attribute == nil {
		return fmt.Errorf("attribute not found")
	}

	attribute.SetValue(value)

	a.Store(attribute)

	return nil
}

func (a *Attributes) Update(entity string, key string, value string, stamp time.Time) error {
	attr := models.Attribute{}
	attr.Entity = entity
	attr.Key = key

	attribute := a.Find(attr.Path())
	if attribute == nil {
		return fmt.Errorf("attribute not found")
	}

	err := attribute.UpdateValue(value, stamp)
	if err != nil {
		return err
	}

	a.Store(attribute)

	return nil
}

func (a *Attributes) Query(entity string, key string) string {
	attr := models.Attribute{}
	attr.Entity = entity
	attr.Key = key

	attribute := a.Find(attr.Path())
	if attribute == nil {
		return ""
	}

	return attribute.Request
}

func (a *Attributes) Compile() []models.Attribute {
	var attributes []models.Attribute
	for _, key := range a.Keys() {
		attribute := a.Find(key)
		if attribute == nil {
			continue
		}
		attributes = append(attributes, *attribute)
	}
	return attributes
}

func (a *Attributes) Find(name string) *models.Attribute {
	res := a.get(name)
	val, ok := res.(*models.Attribute)
	if !ok {
		return nil
	}
	return val
}

func (a *Attributes) Store(attribute *models.Attribute) {
	a.set(attribute.Id, attribute)
	a.emit(attribute.Id, attribute)
}

func LoadAttributes() (m *Attributes) {
	m = &Attributes{}
	m.data = sync.Map{}
	m.Run()
	return m
}
