// Copyright (c) 2022 Braden Nicholson

package domain

import (
	"strconv"
	"time"
	"udap/internal/core/domain/common"
)

type Attribute struct {
	common.Persistent
	Value     string         `json:"value"`
	Updated   time.Time      `json:"updated"`
	Request   string         `json:"request"`
	Requested time.Time      `json:"requested"`
	Entity    string         `json:"entity"`
	Key       string         `json:"key"`
	Type      string         `json:"type"`
	Order     int            `json:"order"`
	Channel   chan Attribute `json:"-" gorm:"-"`
}

const (
	MEDIA  = "media"
	BUFFER = "buffer"
	TOGGLE = "toggle"
	RANGE  = "range"
)

func NewToggleAttribute(entity string) Attribute {
	attribute := Attribute{
		Key:       "on",
		Type:      TOGGLE,
		Entity:    entity,
		Value:     "false",
		Request:   "false",
		Updated:   time.Now(),
		Requested: time.Time{},
		Order:     0,
		Channel:   make(chan Attribute),
	}
	return attribute
}

func NewDimAttribute(entity string) Attribute {
	attribute := Attribute{
		Key:       "dim",
		Type:      RANGE,
		Entity:    entity,
		Value:     "0",
		Request:   "0",
		Updated:   time.Now(),
		Requested: time.Time{},
		Order:     0,
		Channel:   make(chan Attribute),
	}
	return attribute
}

func NewAttribute(key string, variant string, entity string) Attribute {
	attribute := Attribute{
		Key:       key,
		Type:      variant,
		Entity:    entity,
		Value:     "",
		Request:   "",
		Updated:   time.Time{},
		Requested: time.Time{},
		Order:     0,
		Channel:   make(chan Attribute),
	}
	return attribute
}

func (a *Attribute) AsInt() int {
	parsed, err := strconv.ParseInt(a.Value, 10, 0)
	if err != nil {
		return 0
	}
	return int(parsed)
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

type AttributeRepository interface {
	common.Persist[Attribute]
	FindAllByEntity(entity string) (*[]Attribute, error)
	FindByComposite(entity string, key string) (*Attribute, error)
	Register(*Attribute) error
}

type AttributeOperator interface {
	Register(attribute *Attribute) error
	Request(*Attribute, string) error
	Set(*Attribute, string) error
	Update(*Attribute, string, time.Time) error
}

type AttributeService interface {
	Observable
	FindAll() (*[]Attribute, error)
	FindByComposite(entity string, key string) (*Attribute, error)
	FindAllByEntity(entity string) (*[]Attribute, error)
	FindById(id string) (*Attribute, error)
	Create(*Attribute) error
	Register(*Attribute) error
	Request(entity string, key string, value string) error
	Set(entity string, key string, value string) error
	Update(entity string, key string, value string, stamp time.Time) error
	Delete(*Attribute) error
}
