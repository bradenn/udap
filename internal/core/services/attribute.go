// Copyright (c) 2022 Braden Nicholson

package services

import (
	"time"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

func NewAttributeService(repository ports.AttributeRepository, op ports.AttributeOperator) ports.AttributeService {
	return &attributeService{
		repository: repository,
		operator:   op,
	}
}

type attributeService struct {
	repository ports.AttributeRepository
	operator   ports.AttributeOperator
	generic.Watchable[domain.Attribute]
	Logs generic.Watchable[domain.AttributeLog]
}

func (a *attributeService) Summary(key string, start int64, stop int64, window int, mode string) (map[int64]float64, error) {
	return a.repository.Summary(key, start, stop, window, mode)
}

func (a *attributeService) Watch(ref chan<- domain.Mutation) {
	a.Watchable.Watch(ref)
	//a.Logs.Watch(ref)
}

func (a *attributeService) EmitAll() error {
	all, err := a.repository.FindRecent()
	if err != nil {
		return err
	}
	for _, attribute := range *all {
		err = a.Emit(attribute)
		if err != nil {
			return err
		}
	}
	//logs, err := a.FindRecentLogs()
	//if err != nil {
	//	return err
	//}
	//for _, log := range *logs {
	//	err = a.Logs.Emit(log)
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}

func (a *attributeService) FindRecentLogs() (*[]domain.AttributeLog, error) {
	return a.repository.FindRecentLogs()
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

	e.Requested = time.Now()
	e.Request = value
	//log, err := a.repository.Log(&(*e))
	//if err != nil {
	//	return err
	//}

	err = a.operator.Request(e, value)
	if err != nil {
		return err
	}

	e.Value = value

	//e.UpdatedAt = time.Now()
	err = a.repository.Update(e)
	if err != nil {
		return err
	}
	//
	//if log != nil {
	//	err = a.Logs.Emit(*log)
	//	if err != nil {
	//		return err
	//	}
	//}
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
	err = a.repository.Log(e)
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
	attribute.Deleted = true
	_ = a.Emit(*attribute)
	return a.repository.Delete(attribute)
}
