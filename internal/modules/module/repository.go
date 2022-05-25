// Copyright (c) 2022 Braden Nicholson

package module

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

type moduleRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &moduleRepo{
		db: db,
	}
}

func (u moduleRepo) FindAll() ([]*domain.User, error) {
	var target []*domain.User
	if err := u.db.First(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (u moduleRepo) FindById(id string) (*domain.User, error) {
	var target *domain.User
	if err := u.db.Where("id = ?", id).First(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (u moduleRepo) Create(user *domain.User) error {
	if err := u.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (u moduleRepo) FindOrCreate(user *domain.User) error {
	if err := u.db.FirstOrCreate(user).Error; err != nil {
		return err
	}
	return nil
}

func (u moduleRepo) Update(user *domain.User) error {
	if err := u.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (u moduleRepo) Delete(user *domain.User) error {
	if err := u.db.Delete(user).Error; err != nil {
		return err
	}
	return nil
}
