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
		Persistent: common.Persistent{},
		Name:       name,
		Type:       variant,
		Connected:  false,
		Key:        randomSequence(),
	}
}

type EndpointRepository interface {
	common.Persist[Endpoint]
	FindByKey(key string) (*Endpoint, error)
}

type EndpointOperator interface {
	Enroll(*Endpoint, *websocket.Conn) error
	Unenroll(id string) error
	Send(id string, operation string, payload any) error
	SendAll(id string, operation string, payload any) error
	CloseAll() error
}

type EndpointService interface {
	FindAll() (*[]Endpoint, error)
	FindById(id string) (*Endpoint, error)
	FindByKey(key string) (*Endpoint, error)
	Create(*Endpoint) error
	CloseAll() error
	Observable

	Enroll(id string, conn *websocket.Conn) error

	SendAll(target string, operation string, payload any) error
	Send(id string, operation string, payload any) error
	Disconnect(key string) error

	FindOrCreate(*Endpoint) error
	Update(*Endpoint) error
	Delete(*Endpoint) error
}
