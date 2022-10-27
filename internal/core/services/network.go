// Copyright (c) 2022 Braden Nicholson

package services

import (
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

func NewNetworkService(repository ports.NetworkRepository) ports.NetworkService {
	return &networkService{repository: repository}
}

type networkService struct {
	repository ports.NetworkRepository
	generic.Watchable[domain.Network]
}

func (u *networkService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	for _, network := range *all {
		err = u.Emit(network)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *networkService) Register(network *domain.Network) error {
	return u.repository.Register(network)
}

// Repository Mapping

func (u *networkService) FindAll() (*[]domain.Network, error) {
	return u.repository.FindAll()
}

func (u *networkService) FindById(id string) (*domain.Network, error) {
	return u.repository.FindById(id)
}

func (u *networkService) Create(network *domain.Network) error {
	return u.repository.Create(network)
}

func (u *networkService) FindOrCreate(network *domain.Network) error {
	return u.repository.FindOrCreate(network)
}

func (u *networkService) Update(network *domain.Network) error {
	return u.repository.Update(network)
}

func (u *networkService) Delete(network *domain.Network) error {
	return u.repository.Delete(network)
}
