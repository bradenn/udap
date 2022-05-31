// Copyright (c) 2022 Braden Nicholson

package endpoint

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

type endpointRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.EndpointRepository {
	return &endpointRepo{
		db: db,
	}
}

func (u endpointRepo) FindByKey(key string) (*domain.Endpoint, error) {
	var target domain.Endpoint
	if err := u.db.Where("key = ?", key).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (u endpointRepo) FindAll() (*[]domain.Endpoint, error) {
	var target []domain.Endpoint
	if err := u.db.First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (u endpointRepo) FindById(id string) (*domain.Endpoint, error) {
	var target domain.Endpoint
	if err := u.db.Model(&domain.Endpoint{}).Where("id = ?", id).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (u endpointRepo) Create(endpoint *domain.Endpoint) error {
	if err := u.db.Create(endpoint).Error; err != nil {
		return err
	}
	return nil
}

func (u endpointRepo) FindOrCreate(endpoint *domain.Endpoint) error {
	if err := u.db.FirstOrCreate(endpoint).Error; err != nil {
		return err
	}
	return nil
}

func (u endpointRepo) Update(endpoint *domain.Endpoint) error {
	if err := u.db.Save(endpoint).Error; err != nil {
		return err
	}
	return nil
}

func (u endpointRepo) Delete(endpoint *domain.Endpoint) error {
	if err := u.db.Delete(endpoint).Error; err != nil {
		return err
	}
	return nil
}
