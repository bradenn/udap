// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type Device struct {
	common.Persistent
	NetworkId string `json:"networkId" gorm:"-"`
	EntityId  string `json:"entityId" gorm:"-"`
	Name      string `json:"name"`
	Hostname  string `json:"hostname"`
	Mac       string `json:"mac"`
	Ipv4      string `json:"ipv4"`
	Ipv6      string `json:"ipv6"`
}

type DeviceRepository interface {
	common.Persist[Device]
}

type DeviceService interface {
	Observable
	FindAll() (*[]Device, error)
	FindById(id string) (*Device, error)
	Create(*Device) error
	FindOrCreate(*Device) error
	Register(*Device) error
	Update(*Device) error
	Delete(*Device) error
}
