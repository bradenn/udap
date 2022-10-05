// Copyright (c) 2022 Braden Nicholson

package entity

import (
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
	"udap/internal/log"
)

type entityService struct {
	repository ports.EntityRepository
	generic.Watchable[domain.Entity]
}

func (u *entityService) EmitAll() error {
	all, err := u.repository.FindAll()
	if err != nil {
		return err
	}
	for _, et := range *all {
		err = u.Emit(et)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *entityService) mutate(entity *domain.Entity) error {
	err := u.repository.Update(entity)
	if err != nil {
		return err
	}
	err = u.Emit(*entity)
	if err != nil {
		return err
	}
	return nil
}

func (u *entityService) Config(id string, value string) error {
	entity, err := u.FindById(id)
	if err != nil {
		return err
	}
	entity.Config = value
	err = u.Update(entity)
	if err != nil {
		return err
	}
	err = u.Emit(*entity)
	if err != nil {
		return err
	}
	return nil
}

func (u *entityService) Register(entity *domain.Entity) error {
	err := u.repository.Register(entity)
	if err != nil {
		return err
	}
	log.Event("Entity '%s' registered.", entity.Name)
	err = u.Emit(*entity)
	if err != nil {
		return err
	}
	return nil
}

func NewService(repository ports.EntityRepository) ports.EntityService {
	return &entityService{
		repository: repository}
}

func (u entityService) ChangeIcon(id string, icon string) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	byId.Icon = icon
	err = u.mutate(byId)
	if err != nil {
		return err
	}
	return nil
}

// Repository Mapping

func (u entityService) FindByName(name string) (*domain.Entity, error) {
	return u.repository.FindByName(name)
}

func (u *entityService) FindAll() (*[]domain.Entity, error) {
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
