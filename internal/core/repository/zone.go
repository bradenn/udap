// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

type zoneRepo struct {
	generic.Store[domain.Zone]
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) ports.ZoneRepository {
	return &zoneRepo{
		db:    db,
		Store: generic.NewStore[domain.Zone](db),
	}
}

func (z zoneRepo) FindById(id string) (*domain.Zone, error) {
	zone := domain.Zone{}
	err := z.db.Model(&domain.Zone{}).Preload("Entities").Where("id = ?", id).Find(&zone).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func (z zoneRepo) FindByName(name string) (*domain.Zone, error) {
	zone := domain.Zone{}
	err := z.db.Model(&domain.Zone{}).Preload("Entities").Where("name = ?", name).Find(&zone).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func (z zoneRepo) FindAll() (*[]domain.Zone, error) {
	var target []domain.Zone
	if err := z.db.Preload("Entities").Find(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}
