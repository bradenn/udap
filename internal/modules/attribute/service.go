// Copyright (c) 2022 Braden Nicholson

package attribute

import (
	"fmt"
	"time"
	"udap/internal/core/domain"
)

type attributeService struct {
	repository domain.AttributeRepository
	hooks      map[string]chan domain.Attribute
}

func (u attributeService) FindAllByEntity(entity string) (*[]domain.Attribute, error) {
	return u.repository.FindAllByEntity(entity)
}

func NewService(repository domain.AttributeRepository) domain.AttributeService {
	return attributeService{
		repository: repository,
		hooks:      map[string]chan domain.Attribute{},
	}
}

func (u attributeService) Register(attribute *domain.Attribute) error {
	attribute, err := u.repository.FindByComposite(attribute.Entity, attribute.Key)
	if err != nil {
		return err
	}
	if u.hooks[attribute.Id] != nil {
		return fmt.Errorf("attribute already registered")
	}
	u.hooks[attribute.Id] = attribute.Channel
	attribute.Channel = nil
	return nil
}

func (u attributeService) Request(entity string, key string, value string) error {
	return nil
}

func (u attributeService) Set(entity string, key string, value string) error {
	return nil
}

func (u attributeService) Update(entity string, key string, value string, stamp time.Time) error {
	// TODO implement me
	panic("implement me")
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
