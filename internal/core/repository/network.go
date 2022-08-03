// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
)

type networkRepo struct {
	generic.Store[domain.Network]
	db *gorm.DB
}

func NewNetworkRepository(db *gorm.DB) domain.NetworkRepository {
	return &networkRepo{
		db:    db,
		Store: generic.NewStore[domain.Network](db),
	}
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
	var target domain.Network
	if err := u.db.Where("name = ?", name).Find(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}
