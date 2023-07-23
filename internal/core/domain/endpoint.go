// Copyright (c) 2022 Braden Nicholson

package domain

import (
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
	"udap/internal/core/domain/common"
)

type Endpoint struct {
	common.Persistent
	Connection    *websocket.Conn `json:"-" gorm:"-"`
	Name          string          `json:"name" gorm:"unique"`
	Type          string          `json:"type" gorm:"default:'terminal'"`
	Push          string          `json:"push" gorm:"default:'{}'"`
	Notifications bool            `json:"notifications"`
	Connected     bool            `json:"connected"`
	Key           string          `json:"key"`
}

func randomString() string {
	template := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var out string
	rand.Seed(time.Now().Unix())
	for i := 0; i < 4; i++ {
		r := rand.Intn(26)
		u := template[r]
		out += string(u)
	}
	return out
}

func randomNumbers() string {
	template := "0123456789"
	var out string
	rand.Seed(time.Now().Unix())
	for i := 0; i < 4; i++ {
		r := rand.Intn(10)
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
		Key:        fmt.Sprintf("%s%s", randomString(), randomNumbers()),
	}
}
