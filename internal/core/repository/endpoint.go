// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
)

type endpointRepo struct {
	generic.Store[domain.Endpoint]
	db *gorm.DB
}

func NewEndpointRepository(db *gorm.DB) domain.EndpointRepository {
	return &endpointRepo{
		db:    db,
		Store: generic.NewStore[domain.Endpoint](db),
	}
}

func (u endpointRepo) FindByKey(key string) (*domain.Endpoint, error) {
	var target domain.Endpoint
	if err := u.db.Where("key = ?", key).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}
