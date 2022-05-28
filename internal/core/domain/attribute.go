// Copyright (c) 2022 Braden Nicholson

package domain

import "time"

type Attribute struct {
	Persistent
	Value     string    `json:"value"`
	Updated   time.Time `json:"updated"`
	Request   string    `json:"request"`
	Requested time.Time `json:"requested"`
	Entity    string    `json:"entity"`
	Key       string    `json:"key"`
	Type      string    `json:"type"`
	Order     int       `json:"order"`
	Channel   chan Attribute
	// put       FuncPut
	// get       FuncGet
}

type AttributeRepository interface {
	FindAll() (*[]Attribute, error)
	FindAllByEntity(entity string) (*[]Attribute, error)
	FindById(id string) (*Attribute, error)
	FindByComposite(entity string, key string) (*Attribute, error)
	Create(*Attribute) error
	FindOrCreate(*Attribute) error
	Update(*Attribute) error
	Delete(*Attribute) error
}

type AttributeService interface {
	FindAll() (*[]Attribute, error)
	FindAllByEntity(entity string) (*[]Attribute, error)
	FindById(id string) (*Attribute, error)
	Create(*Attribute) error
	Register(*Attribute) error
	Request(entity string, key string, value string) error
	Set(entity string, key string, value string) error
	Update(entity string, key string, value string, stamp time.Time) error
	Delete(*Attribute) error
}
