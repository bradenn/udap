// Copyright (c) 2022 Braden Nicholson

package network

import (
	"udap/internal/core/domain"
)

type networkService struct {
	repository domain.NetworkRepository
}

func (u networkService) Register(network *domain.Network) error {
	return u.repository.Register(network)
}

func NewService(repository domain.NetworkRepository) domain.NetworkService {
	return networkService{repository: repository}
}

// Repository Mapping

func (u networkService) FindAll() ([]*domain.Network, error) {
	return u.repository.FindAll()
}

func (u networkService) FindById(id string) (*domain.Network, error) {
	return u.repository.FindById(id)
}

func (u networkService) Create(network *domain.Network) error {
	return u.repository.Create(network)
}

func (u networkService) FindOrCreate(network *domain.Network) error {
	return u.repository.FindOrCreate(network)
}

func (u networkService) Update(network *domain.Network) error {
	return u.repository.Update(network)
}

func (u networkService) Delete(network *domain.Network) error {
	return u.repository.Delete(network)
}
