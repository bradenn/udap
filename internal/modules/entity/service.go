// Copyright (c) 2022 Braden Nicholson

package entity

import (
	"udap/internal/core/domain"
	"udap/internal/log"
)

type entityService struct {
	repository domain.EntityRepository
}

func (u entityService) Config(id string, value string) error {
	entity, err := u.FindById(id)
	if err != nil {
		return err
	}
	entity.Config = value
	err = u.Update(entity)
	if err != nil {
		return err
	}
	return nil
}

func (u entityService) Register(entity *domain.Entity) error {
	err := u.repository.Register(entity)
	if err != nil {
		return err
	}
	log.Event("Entity '%s' registered.", entity.Name)
	return nil
}

func NewService(repository domain.EntityRepository) domain.EntityService {
	return entityService{repository: repository}
}

// Repository Mapping

func (u entityService) FindAll() (*[]domain.Entity, error) {
	return u.repository.FindAll()
}

func (u entityService) FindById(id string) (*domain.Entity, error) {
	return u.repository.FindById(id)
}

func (u entityService) Create(entity *domain.Entity) error {
	return u.repository.Create(entity)
}

func (u entityService) FindOrCreate(entity *domain.Entity) error {
	return u.repository.FindOrCreate(entity)
}

func (u entityService) Update(entity *domain.Entity) error {
	return u.repository.Update(entity)
}

func (u entityService) Delete(entity *domain.Entity) error {
	return u.repository.Delete(entity)
}
