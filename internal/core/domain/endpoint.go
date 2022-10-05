// Copyright (c) 2022 Braden Nicholson

package domain

import (
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
	"udap/internal/core/domain/common"
)

type Endpoint struct {
	common.Persistent
	Connection *websocket.Conn `json:"-" gorm:"-"`
	Name       string          `json:"name" gorm:"unique"`
	Type       string          `json:"type" gorm:"default:'terminal'"`
	Connected  bool            `json:"connected"`
	Key        string          `json:"key"`
}

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

func NewEndpoint(name string, variant string) *Endpoint {
	return &Endpoint{
		Persistent: common.Persistent{},
		Name:       name,
		Type:       variant,
		Connected:  false,
		Key:        randomSequence(),
	}
}
