// Copyright (c) 2022 Braden Nicholson

package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"udap/internal/cache"
)

type FuncPut func(Attribute) error
type FuncGet func(Attribute) (string, error)

type Attribute struct {
	Updated time.Time `json:"updated"`
	Entity  string    `json:"entity"`
	Key     string    `json:"key"`
	Value   string    `json:"value"`
	Type    string    `json:"type"`
	put     FuncPut
	get     FuncGet
}

func (a *Attribute) Save() error {
	ctx := context.Background()
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return cache.Mem.Set(ctx, strings.ToLower(a.Path()), string(data), 0).Err()
}

func (a *Attribute) Publish() error {
	ctx := context.Background()
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return cache.Mem.Publish(ctx, strings.ToLower(a.Path()), string(data)).Err()
}

func (a *Attribute) Retrieve() error {
	ctx := context.Background()
	result, err := cache.Mem.Get(ctx, strings.ToLower(a.Path())).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), a)
}

func (a *Attribute) Path() string {
	return fmt.Sprintf("entity.%s.attribute.%s", a.Entity, a.Key)
}

func (a *Attribute) Request(val string) error {

	if a.put == nil {
		return fmt.Errorf("attribute put function not connected")
	}
	a.Value = val
	err := a.put(*a)
	if err != nil {
		return err
	}

	a.Updated = time.Now()
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

func NewMediaEntity(name string, module string) *Entity {
	e := Entity{
		Name:   name,
		Type:   "media",
		Module: module,
	}
	return &e
}

func NewEntity(name string, module string) *Entity {
	e := Entity{
		Name:   name,
		Type:   "switch",
		Module: module,
	}
	return &e
}

func NewSpectrum(name string, module string) *Entity {

	e := Entity{
		Name:   name,
		Type:   "spectrum",
		Module: module,
	}
	return &e
}

func NewSwitch(name string, module string) *Entity {

	e := Entity{
		Name:   name,
		Type:   "switch",
		Module: module,
	}
	return &e
}
