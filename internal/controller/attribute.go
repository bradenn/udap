// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"encoding/json"
	"math/rand"
	"sync"
	"time"
	"udap/internal/bond"
	"udap/internal/models"
)

type Attributes struct {
	PolyBuffer
}

// randomSequence generates a random id for use as a key
func randomSequence() string {
	template := "abcdefghijklmnopqrstuvwxyz"
	var out string
	rand.Seed(time.Now().Unix())
	for i := 0; i < 8; i++ {
		r := rand.Intn(26)
		u := template[r]
		out += string(u)
	}
	return out
}

func (a *Attributes) Handle(event bond.Msg) (res any, err error) {
	switch event.Operation {
	case "request":
		return a.request(event)
	}
	return nil, nil
}

func (a *Attributes) Register(attribute *models.Attribute) error {
	attribute.Id = attribute.Path()
	a.Store(attribute)
	return nil
}

func (a *Attributes) request(event bond.Msg) (any, error) {
	attr := models.Attribute{}
	err := json.Unmarshal([]byte(event.Payload), &attr)
	if err != nil {

		return nil, err
	}

	attribute := a.Find(attr.Id)
	if err != nil {
		return nil, err
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

func LoadAttributes() (m *Attributes) {
	m = &Attributes{}
	m.data = sync.Map{}
	return m
}
