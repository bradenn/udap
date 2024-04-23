// Copyright (c) 2022 Braden Nicholson

package domain

import (
	"strconv"
	"time"
	"udap/internal/core/domain/common"
)

type AttributeLog struct {
	common.Persistent
	Attribute string    `json:"attribute"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Time      time.Time `json:"time"`
}

type Attribute struct {
	common.Persistent
	Value     string         `json:"value"`
	Updated   time.Time      `json:"lastUpdated"`
	Request   string         `json:"request"`
	Requested time.Time      `json:"requested"`
	Entity    string         `json:"entity"`
	Serial    string         `json:"serial"`
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
	parsed, err := strconv.ParseInt(a.Value, 10, 32)
	if err != nil {
		return 0
	}
	return int(parsed)
}

func (a *Attribute) ToLog() AttributeLog {
	return AttributeLog{
		Attribute: a.Id,
		To:        a.Request,
		From:      a.Value,
		Time:      time.Now(),
	}
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
