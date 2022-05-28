// Copyright (c) 2022 Braden Nicholson

package entity

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

type entityRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.EntityRepository {
	return &entityRepo{
		db: db,
	}
}

func (u entityRepo) FindAll() (*[]domain.Entity, error) {
	var target *[]domain.Entity
	if err := u.db.First(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (u entityRepo) FindById(id string) (*domain.Entity, error) {
	var target *domain.Entity
	if err := u.db.Where("id = ?", id).First(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (u entityRepo) Create(entity *domain.Entity) error {
	if err := u.db.Create(entity).Error; err != nil {
		return err
	}
	return nil
}

func (u entityRepo) FindOrCreate(entity *domain.Entity) error {
	if err := u.db.FirstOrCreate(entity).Error; err != nil {
		return err
	}
	return nil
}

func (u entityRepo) Update(entity *domain.Entity) error {
	if err := u.db.Save(entity).Error; err != nil {
		return err
	}
	return nil
}

func (u entityRepo) Delete(entity *domain.Entity) error {
	if err := u.db.Delete(entity).Error; err != nil {
		return err
	}
	return nil
}
