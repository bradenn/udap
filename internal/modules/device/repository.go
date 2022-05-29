// Copyright (c) 2022 Braden Nicholson

package device

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

type deviceRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.DeviceRepository {
	return &deviceRepo{
		db: db,
	}
}

func (u deviceRepo) FindAll() ([]*domain.Device, error) {
	var target []*domain.Device
	if err := u.db.First(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (u deviceRepo) FindById(id string) (*domain.Device, error) {
	var target *domain.Device
	if err := u.db.Where("id = ?", id).First(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (u deviceRepo) Create(device *domain.Device) error {
	if err := u.db.Create(device).Error; err != nil {
		return err
	}
	return nil
}

func (u deviceRepo) FindOrCreate(device *domain.Device) error {
	if err := u.db.FirstOrCreate(device).Error; err != nil {
		return err
	}
	return nil
}

func (u deviceRepo) Update(device *domain.Device) error {
	if err := u.db.Save(device).Error; err != nil {
		return err
	}
	return nil
}

func (u deviceRepo) Delete(device *domain.Device) error {
	if err := u.db.Delete(device).Error; err != nil {
		return err
	}
	return nil
}
