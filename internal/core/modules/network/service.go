// Copyright (c) 2022 Braden Nicholson

package network

import (
	"fmt"
	"udap/internal/core/domain"
)

type networkService struct {
	repository domain.NetworkRepository
	channel    chan<- domain.Mutation
}

func (u *networkService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	for _, network := range *all {
		err = u.emit(&network)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *networkService) emit(network *domain.Network) error {
	if u.channel == nil {
		return nil
	}
	u.channel <- domain.Mutation{
		Status:    "update",
		Operation: "network",
		Body:      *network,
		Id:        network.Id,
	}
	return nil
}

func (u *networkService) Watch(mut chan<- domain.Mutation) error {
	if u.channel != nil {
		return fmt.Errorf("channel already set")
	}
	u.channel = mut

	return nil
}

func (u networkService) Register(network *domain.Network) error {
	return u.repository.Register(network)
}

func NewService(repository domain.NetworkRepository) domain.NetworkService {
	return &networkService{repository: repository}
}

// Repository Mapping

func (u networkService) FindAll() (*[]domain.Network, error) {
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
