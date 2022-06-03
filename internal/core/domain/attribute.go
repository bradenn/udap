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
	// put       FuncPut
	// get       FuncGet
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
	FindAllByEntity(entity string) (*[]Attribute, error)
	FindById(id string) (*Attribute, error)
	Create(*Attribute) error
	Register(*Attribute) error
	Request(entity string, key string, value string) error
	Set(entity string, key string, value string) error
	Update(entity string, key string, value string, stamp time.Time) error
	Delete(*Attribute) error
}
