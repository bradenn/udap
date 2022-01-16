// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"encoding/json"
	"sync"
	"udap/internal/bond"
	"udap/internal/models"
)

type Attributes struct {
	PolyBuffer
}

func (a *Attributes) Handle(event bond.Msg) (res any, err error) {
	switch event.Operation {
	case "request":
		return a.request(event)
	}
	return nil, nil
}

func (a *Attributes) Register(attribute *models.Attribute) error {
	err := attribute.Save()
	if err != nil {
		return err
	}
	a.set(attribute.Path(), attribute)
	return nil
}

func (a *Attributes) request(event bond.Msg) (any, error) {
	attr := &models.Attribute{}

	err := json.Unmarshal([]byte(event.Payload), attr)
	if err != nil {
		return nil, err
	}
	value := attr.Value

	attr = a.Find(attr.Path())
	attr.Value = value
	a.set(attr.Path(), attr)
	err = attr.Request(value)
	if err != nil {
		return nil, err
	}

	err = attr.Save()
	if err != nil {
		return nil, err
	}

	err = attr.Publish()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (a *Attributes) Put(entity string, key string, value string) error {
	attr := &models.Attribute{}
	attr.Entity = entity
	attr.Key = key
	attr = a.Find(attr.Path())
	attr.Value = value

	err := attr.Save()
	if err != nil {
		return err
	}

	err = attr.Publish()
	if err != nil {
		return err
	}
	a.set(attr.Path(), attr)
	return nil
}

func (a *Attributes) Compile() []models.Attribute {
	var attributes []models.Attribute
	for _, key := range a.Keys() {
		attribute := a.Find(key)
		attributes = append(attributes, *attribute)
	}
	return attributes
}

func (a *Attributes) Find(name string) *models.Attribute {
	return a.get(name).(*models.Attribute)
}

func LoadAttributes() (m *Attributes) {
	m = &Attributes{}
	m.data = sync.Map{}
	return m
}
