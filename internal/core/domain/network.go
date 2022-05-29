// Copyright (c) 2022 Braden Nicholson

package domain

type Network struct {
	Persistent
	Name   string `json:"name"`
	Dns    string `json:"dns"`
	Router string `json:"index"`
	Lease  string `json:"lease"`
	Mask   string `json:"mask"`
	Range  string `json:"range"`
}

type NetworkRepository interface {
	FindAll() ([]*Network, error)
	FindById(id string) (*Network, error)
	FindByName(name string) (*Network, error)
	Create(*Network) error
	Register(*Network) error
	FindOrCreate(*Network) error
	Update(*Network) error
	Delete(*Network) error
}

type NetworkService interface {
	FindAll() ([]*Network, error)
	FindById(id string) (*Network, error)
	Create(*Network) error
	FindOrCreate(*Network) error
	Register(*Network) error
	Update(*Network) error
	Delete(*Network) error
}
