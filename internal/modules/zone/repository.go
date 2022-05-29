// Copyright (c) 2022 Braden Nicholson

package zone

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

type zoneRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.ZoneRepository {
	return &zoneRepo{
		db: db,
	}
}

func (u zoneRepo) FindAll() ([]*domain.Zone, error) {
	var target []*domain.Zone
	if err := u.db.First(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (u zoneRepo) FindById(id string) (*domain.Zone, error) {
	var target *domain.Zone
	if err := u.db.Where("id = ?", id).First(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (u zoneRepo) Create(zone *domain.Zone) error {
	if err := u.db.Create(zone).Error; err != nil {
		return err
	}
	return nil
}

func (u zoneRepo) FindOrCreate(zone *domain.Zone) error {
	if err := u.db.FirstOrCreate(zone).Error; err != nil {
		return err
	}
	return nil
}

func (u zoneRepo) Update(zone *domain.Zone) error {
	if err := u.db.Save(zone).Error; err != nil {
		return err
	}
	return nil
}

func (u zoneRepo) Delete(zone *domain.Zone) error {
	if err := u.db.Delete(zone).Error; err != nil {
		return err
	}
	return nil
}
