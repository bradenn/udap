// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"time"
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type DeviceRepository interface {
	common.Persist[domain.Device]
}

type DeviceService interface {
	domain.Observable
	FindAll() (*[]domain.Device, error)
	FindById(id string) (*domain.Device, error)
	Create(*domain.Device) error
	FindOrCreate(*domain.Device) error
	Register(*domain.Device) error
	Update(*domain.Device) error
	Delete(*domain.Device) error
	Ping(id string, latency time.Duration) error
	Utilization(id string, utilization domain.Utilization) error
}
