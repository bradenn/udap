// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"udap/internal/bond"
	"udap/internal/models"
)

type Attributes struct {
	PolyBuffer
}

type IAttribute interface {
	GetValue(value string)
	GetRequest(value string)
}

type Observer chan bytes.Buffer

func (e *Endpoints) Listen() {
	c := make(Observer, 2)
	for {

		select {
		case <-c: // Called when the attribute state changes

		case <-time.After(time.Second * 1): // Called if not updated by reaction for more than 5 seconds

		}
	}
}

func (a *Attributes) Handle(event bond.Msg) (res any, err error) {
	switch event.Operation {
	case "request":
		return a.request(event)
	case "poll":
		return a.poll(event)
	default:
		return nil, fmt.Errorf("operaiton not found")
	}
}

func (a *Attributes) Register(attribute *models.Attribute) error {
	attribute.Id = attribute.Path()
	a.Store(attribute)
	return nil
}

func (a *Attributes) poll(event bond.Msg) (any, error) {
	attr := models.Attribute{}

	err := json.Unmarshal([]byte(event.Payload), &attr)
	if err != nil {
		return nil, err
	}
	attribute := a.Find(attr.Id)
	if attribute == nil {
		return nil, fmt.Errorf("attribute '%s' not found", attr.Id)
	}

	ok, err := attribute.Poll()
	if err != nil {
		return nil, err
	}

	if ok {
		a.Store(attribute)
	}

	return nil, nil
}

func (a *Attributes) request(event bond.Msg) (any, error) {
	attr := models.Attribute{}

	err := json.Unmarshal([]byte(event.Payload), &attr)
	if err != nil {
		return nil, err
	}
	attribute := a.Find(attr.Id)
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

func (a *Attributes) Update(entity string, key string, value string) error {
	attr := models.Attribute{}
	attr.Entity = entity
	attr.Key = key

	attribute := a.Find(attr.Path())
	err := attribute.UpdateValue(value)
	if err != nil {
		return err
	}

	a.Store(attribute)

	return nil
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
	err := attribute.CacheIn()
	if err != nil {
		return
	}
	a.set(attribute.Id, attribute)
}

func (a *Attributes) Observe(func(attribute models.Attribute) error) {

}

func LoadAttributes() (m *Attributes) {
	m = &Attributes{}
	m.data = sync.Map{}
	return m
}
