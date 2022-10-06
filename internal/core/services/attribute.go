// Copyright (c) 2022 Braden Nicholson

package services

import (
	"gorm.io/gorm"
	"time"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
	"udap/internal/core/repository"
)

func NewAttributeService(db *gorm.DB, op ports.AttributeOperator) ports.AttributeService {
	repo := repository.NewAttributeRepository(db)
	return &attributeService{
		repository: repo,
		operator:   op,
	}
}

type attributeService struct {
	repository ports.AttributeRepository
	operator   ports.AttributeOperator
	generic.Watchable[domain.Attribute]
}

func (a *attributeService) EmitAll() error {
	all, err := a.FindAll()
	if err != nil {
		return err
	}
	for _, attribute := range *all {
		err = a.Emit(attribute)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *attributeService) FindAllByEntity(entity string) (*[]domain.Attribute, error) {
	return a.repository.FindAllByEntity(entity)
}

func (a *attributeService) Register(attribute *domain.Attribute) error {
	err := a.repository.Register(attribute)
	if err != nil {
		return err
	}
	err = a.operator.Register(attribute)
	if err != nil {
		return err
	}

	err = a.Emit(*attribute)
	if err != nil {
		return err
	}
	return nil
}

func (a *attributeService) Request(entity string, key string, value string) error {
	e, err := a.repository.FindByComposite(entity, key)
	if err != nil {
		return err
	}
	err = a.operator.Request(e, value)
	if err != nil {
		return err
	}
	err = a.repository.Update(e)
	if err != nil {
		return err
	}
	err = a.Emit(*e)
	if err != nil {
		return err
	}
	return nil
}

func (a *attributeService) Set(entity string, key string, value string) error {
	e, err := a.repository.FindByComposite(entity, key)
	if err != nil {
		return err
	}
	err = a.operator.Update(e, value, time.Now())
	if err != nil {
		return err
	}
	err = a.repository.Update(e)
	if err != nil {
		return err
	}
	err = a.Emit(*e)
	if err != nil {
		return err
	}
	return nil
}

func (a *attributeService) Update(entity string, key string, value string, stamp time.Time) error {
	e, err := a.repository.FindByComposite(entity, key)
	if err != nil {
		return err
	}
	err = a.operator.Update(e, value, stamp)
	if err != nil {
		return err
	}
	err = a.repository.Update(e)
	if err != nil {
		return err
	}
	err = a.Emit(*e)
	if err != nil {
		return err
	}
	return nil
}

func (a *attributeService) FindByComposite(entity string, key string) (*domain.Attribute, error) {
	return a.repository.FindByComposite(entity, key)
}

// Repository Mapping

func (a *attributeService) FindAll() (*[]domain.Attribute, error) {
	return a.repository.FindAll()
}

func (a *attributeService) FindById(id string) (*domain.Attribute, error) {
	return a.repository.FindById(id)
}

func (a *attributeService) Create(attribute *domain.Attribute) error {
	return a.repository.Create(attribute)
}

func (a *attributeService) FindOrCreate(attribute *domain.Attribute) error {
	return a.repository.FindOrCreate(attribute)
}

func (a *attributeService) Delete(attribute *domain.Attribute) error {
	return a.repository.Delete(attribute)
}