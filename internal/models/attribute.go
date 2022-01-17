// Copyright (c) 2022 Braden Nicholson

package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"udap/internal/cache"
)

type FuncPut func(value string) error
type FuncGet func() (string, error)

type Attribute struct {
	Id      string    `json:"id"`
	Value   string    `json:"value"`
	Updated time.Time `json:"updated"`

	Request   string    `json:"request"`
	Requested time.Time `json:"requested"`

	Entity string `json:"entity"`
	Key    string `json:"key"`
	Type   string `json:"type"`
	Order  int    `json:"order"`
	put    FuncPut
	get    FuncGet
}

func (a *Attribute) Path() string {
	return fmt.Sprintf("%s.%s", a.Entity, a.Key)
}

func (a *Attribute) CacheIn() error {
	marshal, err := json.Marshal(a)
	if err != nil {
		return err
	}
	err = cache.PutLn(string(marshal), a.Id)
	if err != nil {
		return err
	}
	return nil
}

func (a *Attribute) CacheOut() error {

	ln, err := cache.GetLn(a.Entity, a.Key)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(ln.(string)), a)
	if err != nil {
		return err
	}
	return nil
}

func (a *Attribute) UpdateValue(val string) error {
	a.Value = val
	a.Updated = time.Now()
	return nil
}

func (a *Attribute) SendRequest(val string) error {
	if a.put == nil {
		return fmt.Errorf("attribute put function not connected")
	}

	a.Request = val
	a.Requested = time.Now()

	err := a.put(val)
	if err != nil {
		return err
	}
	a.Value = val
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
	parsed, err := strconv.ParseInt(a.Value, 10, 64)
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
