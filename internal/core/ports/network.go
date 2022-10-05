// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type NetworkRepository interface {
	common.Persist[domain.Network]
	FindByName(name string) (*domain.Network, error)
	Register(*domain.Network) error
}

type NetworkService interface {
	domain.Observable
	FindAll() (*[]domain.Network, error)
	FindById(id string) (*domain.Network, error)
	Create(*domain.Network) error
	FindOrCreate(*domain.Network) error
	Register(*domain.Network) error
	Update(*domain.Network) error
	Delete(*domain.Network) error
}
