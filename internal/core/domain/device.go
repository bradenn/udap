// Copyright (c) 2022 Braden Nicholson

package domain

type Device struct {
	Persistent
	NetworkId string `json:"networkId" gorm:"-"`
	EntityId  string `json:"entityId" gorm:"-"`
	Name      string `json:"name"`
	Hostname  string `json:"hostname"`
	Mac       string `json:"mac"`
	Ipv4      string `json:"ipv4"`
	Ipv6      string `json:"ipv6"`
}

type DeviceRepository interface {
	FindAll() ([]*Device, error)
	FindById(id string) (*Device, error)
	Create(*Device) error
	FindOrCreate(*Device) error
	Update(*Device) error
	Delete(*Device) error
}

type DeviceService interface {
	FindAll() ([]*Device, error)
	FindById(id string) (*Device, error)
	Create(*Device) error
	FindOrCreate(*Device) error
	Register(*Device) error
	Update(*Device) error
	Delete(*Device) error
}
