// Copyright (c) 2022 Braden Nicholson

package entity

import (
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/log"
)

type entityService struct {
	repository domain.EntityRepository
	channel    chan<- domain.Mutation
	generic.Watchable
}

func (u entityService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	for _, entity := range *all {
		err = u.emit(entity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u entityService) mutate(entity *domain.Entity) error {
	err := u.repository.Update(entity)
	if err != nil {
		return err
	}
	err = u.emit(*entity)
	if err != nil {
		return err
	}
	return nil
}

func (u entityService) emit(entity domain.Entity) error {
	err := u.Emit(entity, entity.Id)
	if err != nil {
		log.Err(err)
	}
	return nil
}

// func (u *entityService) Watch(mut chan<- domain.Mutation) error {
// 	if u.channel != nil {
// 		return fmt.Errorf("channel already set")
// 	}
// 	u.channel = mut
//
// 	return nil
// }

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
	err = u.emit(*entity)
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
	err = u.emit(*entity)
	if err != nil {
		return err
	}
	return nil
}

func NewService(repository domain.EntityRepository) domain.EntityService {
	return &entityService{
		Watchable:  generic.NewWatchable("entity"),
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
