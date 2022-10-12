// Copyright (c) 2022 Braden Nicholson

package services

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
	"udap/internal/core/repository"
)

func NewNetworkService(db *gorm.DB) ports.NetworkService {
	repo := repository.NewNetworkRepository(db)
	return &networkService{repository: repo}
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
