// Copyright (c) 2022 Braden Nicholson

package attribute

import (
	"time"
	"udap/internal/core/domain"
)

type attributeService struct {
	repository domain.AttributeRepository
	operator   domain.AttributeOperator
	channel    chan<- domain.Attribute
}

func (u *attributeService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	attributes := *all
	for _, attribute := range attributes {
		err = u.push(&attribute)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *attributeService) push(attribute *domain.Attribute) error {
	u.channel <- *attribute
	return nil
}

func (u *attributeService) Watch(channel chan<- domain.Attribute) error {
	u.channel = channel
	return nil
}

func (u attributeService) FindAllByEntity(entity string) (*[]domain.Attribute, error) {
	return u.repository.FindAllByEntity(entity)
}

func NewService(repository domain.AttributeRepository, operator domain.AttributeOperator) domain.AttributeService {
	return &attributeService{
		repository: repository,
		operator:   operator,
		channel:    nil,
	}
}

func (u attributeService) Register(attribute *domain.Attribute) error {
	err := u.repository.Register(attribute)
	if err != nil {
		return err
	}
	err = u.operator.Register(attribute)
	if err != nil {
		return err
	}

	return nil
}

func (u attributeService) Request(entity string, key string, value string) error {
	e, err := u.repository.FindByComposite(entity, key)
	if err != nil {
		return err
	}
	err = u.operator.Request(e, value)
	if err != nil {
		return err
	}
	err = u.push(e)
	if err != nil {
		return err
	}
	return nil
}

func (u attributeService) Set(entity string, key string, value string) error {
	e, err := u.repository.FindByComposite(entity, key)
	if err != nil {
		return err
	}
	err = u.operator.Set(e, value)
	if err != nil {
		return err
	}
	err = u.push(e)
	if err != nil {
		return err
	}
	return nil
}

func (u attributeService) Update(entity string, key string, value string, stamp time.Time) error {
	e, err := u.repository.FindByComposite(entity, key)
	if err != nil {
		return err
	}
	err = u.operator.Update(e, value, stamp)
	if err != nil {
		return err
	}
	err = u.push(e)
	if err != nil {
		return err
	}
	return nil
}

// Repository Mapping

func (u attributeService) FindAll() (*[]domain.Attribute, error) {
	return u.repository.FindAll()
}

func (u attributeService) FindById(id string) (*domain.Attribute, error) {
	return u.repository.FindById(id)
}

func (u attributeService) Create(attribute *domain.Attribute) error {
	return u.repository.Create(attribute)
}

func (u attributeService) FindOrCreate(attribute *domain.Attribute) error {
	return u.repository.FindOrCreate(attribute)
}

func (u attributeService) Delete(attribute *domain.Attribute) error {
	return u.repository.Delete(attribute)
}
