// Copyright (c) 2022 Braden Nicholson

package network

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

type networkRepo struct {
	db *gorm.DB
}

func (u networkRepo) Register(network *domain.Network) error {
	if network.Id == "" {
		err := u.db.Model(&domain.Network{}).Where("name = ?", network.Name).FirstOrCreate(network).Error
		if err != nil {
			return err
		}
	} else {
		err := u.db.Model(&domain.Network{}).Where("id = ?", network.Id).First(network).Error
		if err != nil {
			return err
		}
	}
	err := u.db.Model(&domain.Network{}).Where("id = ?", network.Id).Save(network).Error
	if err != nil {
		return err
	}
	return nil
}

func (u networkRepo) FindByName(name string) (*domain.Network, error) {
	var target *domain.Network
	if err := u.db.Where("name = ?", name).Find(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func NewRepository(db *gorm.DB) domain.NetworkRepository {
	return &networkRepo{
		db: db,
	}
}

func (u networkRepo) FindAll() ([]*domain.Network, error) {
	var target []*domain.Network
	if err := u.db.First(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (u networkRepo) FindById(id string) (*domain.Network, error) {
	var target *domain.Network
	if err := u.db.Where("id = ?", id).First(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (u networkRepo) Create(network *domain.Network) error {
	if err := u.db.Create(network).Error; err != nil {
		return err
	}
	return nil
}

func (u networkRepo) FindOrCreate(network *domain.Network) error {
	if err := u.db.FirstOrCreate(network).Error; err != nil {
		return err
	}
	return nil
}

func (u networkRepo) Update(network *domain.Network) error {
	if err := u.db.Save(network).Error; err != nil {
		return err
	}
	return nil
}

func (u networkRepo) Delete(network *domain.Network) error {
	if err := u.db.Delete(network).Error; err != nil {
		return err
	}
	return nil
}
