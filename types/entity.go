// Copyright (c) 2021 Braden Nicholson

package types

import (
	"github.com/google/uuid"
)

const (
	SWITCH = "switch"
	RGB    = "rgb"
	RGBW   = "rgbw"
	RGBCCT = "rgbcct"
)

type Entity struct {
	Persistent
	InstanceId uuid.UUID `json:"instanceId"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	State      string    `json:"state"`
}

func NewEntity(name string, entityType string, instanceId string) (err error) {
	return nil
}

func (e *Entity) SetState(state string) (err error) {

	return nil
}

type Device struct {
	Persistent
	InstanceId uuid.UUID `json:"instanceId"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
}
