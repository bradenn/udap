// Copyright (c) 2021 Braden Nicholson

package models

import (
	"udap/internal/store"
)

type Log struct {
	store.Persistent
	EntityId string `json:"entityId"`
	Power    string `json:"power"`
	Mode     string `json:"mode"`
	Level    int    `json:"level"`
	CCT      int    `json:"cct"`
}
