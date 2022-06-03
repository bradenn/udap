// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type Network struct {
	common.Persistent
	Name   string `json:"name"`
	Dns    string `json:"dns"`
	Router string `json:"index"`
	Lease  string `json:"lease"`
	Mask   string `json:"mask"`
	Range  string `json:"range"`
}

type NetworkRepository interface {
	common.Persist[Network]
	FindByName(name string) (*Network, error)
	Register(*Network) error
}

type NetworkService interface {
	Observable
	FindAll() (*[]Network, error)
	FindById(id string) (*Network, error)
	Create(*Network) error
	FindOrCreate(*Network) error
	Register(*Network) error
	Update(*Network) error
	Delete(*Network) error
}
