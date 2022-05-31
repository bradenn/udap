// Copyright (c) 2022 Braden Nicholson

package user

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

type userRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.UserRepository {
	return &userRepo{
		db: db,
	}
}

func (u userRepo) FindAll() (*[]domain.User, error) {
	var target []domain.User
	if err := u.db.Find(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (u userRepo) FindById(id string) (*domain.User, error) {
	var target *domain.User
	if err := u.db.Where("id = ?", id).First(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (u userRepo) Create(user *domain.User) error {
	if err := u.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (u userRepo) FindOrCreate(user *domain.User) error {
	if err := u.db.FirstOrCreate(user).Error; err != nil {
		return err
	}
	return nil
}

func (u userRepo) Update(user *domain.User) error {
	if err := u.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (u userRepo) Delete(user *domain.User) error {
	if err := u.db.Delete(user).Error; err != nil {
		return err
	}
	return nil
}
