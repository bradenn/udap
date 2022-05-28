// Copyright (c) 2022 Braden Nicholson

package domain

type Endpoint struct {
	Persistent
	Name      string `json:"name" gorm:"unique"`
	Type      string `json:"type" gorm:"default:'terminal'"`
	Connected bool   `json:"connected"`
	Key       string `json:"key"`
}

type EndpointRepository interface {
	FindAll() ([]*Endpoint, error)
	FindById(id string) (*Endpoint, error)
	FindByKey(key string) (*Endpoint, error)
	Create(*Endpoint) error
	FindOrCreate(*Endpoint) error
	Update(*Endpoint) error
	Delete(*Endpoint) error
}

type EndpointService interface {
	FindAll() ([]*Endpoint, error)
	FindById(id string) (*Endpoint, error)
	FindByKey(key string) (*Endpoint, error)
	Create(*Endpoint) error

	Enroll(key string) (*Endpoint, error)
	Disconnect(key string) error

	FindOrCreate(*Endpoint) error
	Update(*Endpoint) error
	Delete(*Endpoint) error
}
