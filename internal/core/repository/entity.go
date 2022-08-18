// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
)

type entityRepo struct {
	generic.Store[domain.Entity]
	db *gorm.DB
}

func (u entityRepo) FindAll() (*[]domain.Entity, error) {
	var res []domain.Entity
	if err := u.db.Model(domain.Entity{}).Find(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

func (u entityRepo) FindByName(name string) (*domain.Entity, error) {
	entity := domain.Entity{}
	err := u.db.Model(&domain.Entity{}).Where("name = ?", name).Find(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (u entityRepo) Register(e *domain.Entity) error {
	if e.Id == "" {
		err := u.db.Model(&domain.Entity{}).Where("name = ? AND module = ?", e.Name, e.Module).FirstOrCreate(e).Error
		if err != nil {
			return err
		}
	} else {
		err := u.db.Model(&domain.Entity{}).Where("name = ?", e.Name).First(e).Error
		if err != nil {
			return err
		}
	}
	err := u.db.Model(&domain.Entity{}).Where("name = ?", e.Name).Save(e).Error
	if err != nil {
		return err
	}
	return nil
}

func NewEntityRepository(db *gorm.DB) domain.EntityRepository {
	return &entityRepo{
		db:    db,
		Store: generic.NewStore[domain.Entity](db),
	}
}
