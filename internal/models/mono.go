// Copyright (c) 2021 Braden Nicholson

package models

import (
	"encoding/json"
)

type Mono struct {
	Value float32 `json:"value"`
}

func (m *Mono) Marshal() []byte {
	marshal, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return marshal
}

func (m *Mono) Unmarshal(data []byte) {
	err := json.Unmarshal(data, m)
	if err != nil {
		return
	}
}

func NewDimmer(name string, module string) *Entity {
	mono := Mono{
		Value: 0,
	}
	e := Entity{
		Name:   name,
		Type:   "dimmer",
		Module: module,
		State:  mono.Marshal(),
	}
	return &e
}

type Rgbw struct {
	R   int `json:"r"`
	G   int `json:"g"`
	B   int `json:"b"`
	CCT int `json:"cct"`
}

func (r *Rgbw) Unmarshal(data []byte) {
	err := json.Unmarshal(data, r)
	if err != nil {
		return
	}
}

func (r *Rgbw) Marshal() []byte {
	marshal, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return marshal
}

func NewWildcardDevice(name string, module string, state State) *Entity {

	e := Entity{
		Name:   name,
		Type:   "wildcard",
		Module: module,
		State:  json.RawMessage(state),
	}
	return &e
}

func NewSwitch(name string, module string) *Entity {
	mono := Mono{
		Value: 0,
	}
	e := Entity{
		Name:   name,
		Type:   "switch",
		Module: module,
		State:  mono.Marshal(),
	}
	return &e
}
