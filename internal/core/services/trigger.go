// Copyright (c) 2022 Braden Nicholson

package services

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
	"udap/internal/core/repository"
)

func NewTriggerService(db *gorm.DB) ports.TriggerService {
	repo := repository.NewTriggerRepository(db)
	return &triggerService{repository: repo}
}

type triggerService struct {
	repository ports.UserRepository
	generic.Watchable[domain.User]
}

func (u *triggerService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	for _, user := range *all {
		err = u.Emit(user)
		if err != nil {
			return err
		}
	}
	return nil
}

// Repository Mapping

func (u *triggerService) FindAll() (*[]domain.User, error) {
	return u.repository.FindAll()
}

func (u *triggerService) FindById(id string) (*domain.User, error) {
	return u.repository.FindById(id)
}

func (u *triggerService) Create(user *domain.User) error {
	return u.repository.Create(user)
}

func (u *triggerService) FindOrCreate(user *domain.User) error {
	return u.repository.FindOrCreate(user)
}

func (u *triggerService) Update(user *domain.User) error {
	return u.repository.Update(user)
}

func (u *triggerService) Delete(user *domain.User) error {
	return u.repository.Delete(user)
}
