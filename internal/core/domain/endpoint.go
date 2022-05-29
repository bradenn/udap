// Copyright (c) 2022 Braden Nicholson

package domain

import (
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
)

type Endpoint struct {
	Persistent
	Name      string `json:"name" gorm:"unique"`
	Type      string `json:"type" gorm:"default:'terminal'"`
	Connected bool   `json:"connected"`
	Key       string `json:"key"`
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
		Persistent: Persistent{},
		Name:       name,
		Type:       variant,
		Connected:  false,
		Key:        randomSequence(),
	}
}

type EndpointRepository interface {
	FindAll() (*[]Endpoint, error)
	FindById(id string) (*Endpoint, error)
	FindByKey(key string) (*Endpoint, error)
	Create(*Endpoint) error
	FindOrCreate(*Endpoint) error
	Update(*Endpoint) error
	Delete(*Endpoint) error
}

type EndpointOperator interface {
	Enroll(*Endpoint, *websocket.Conn) error
	Send(id string, operation string, payload any) error
}

type EndpointService interface {
	FindAll() (*[]Endpoint, error)
	FindById(id string) (*Endpoint, error)
	FindByKey(key string) (*Endpoint, error)
	Create(*Endpoint) error

	Enroll(id string, conn *websocket.Conn) error
	Send(id string, operation string, payload any) error
	Disconnect(key string) error

	FindOrCreate(*Endpoint) error
	Update(*Endpoint) error
	Delete(*Endpoint) error
}
