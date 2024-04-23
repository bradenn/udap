// Copyright (c) 2022 Braden Nicholson

package services

import (
	"time"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

func NewTriggerService(repository ports.TriggerRepository, operator ports.TriggerOperator) ports.TriggerService {
	return &triggerService{repository: repository, operator: operator}
}

type triggerService struct {
	repository ports.TriggerRepository
	operator   ports.TriggerOperator
	generic.Watchable[domain.Trigger]
}

func (u *triggerService) Register(trigger *domain.Trigger) error {
	err := u.repository.Register(trigger)
	if err != nil {
		return err
	}
	return nil
}

func (u *triggerService) Trigger(name string) error {
	trigger, err := u.repository.FindByName(name)
	if err != nil {
		return err
	}
	err = u.operator.Run(*trigger)
	if err != nil {
		return err
	}
	trigger.LastTrigger = time.Now()
	err = u.Update(trigger)
	if err != nil {
		return err
	}
	err = u.Emit(*trigger)
	if err != nil {
		return err
	}
	return nil
}

func (u *triggerService) TriggerCustom(name string, key string, value string) error {
	trigger, err := u.repository.FindByName(name)
	if err != nil {
		return err
	}
	err = u.operator.RunCustom(*trigger, key, value)
	if err != nil {
		return err
	}
	trigger.LastTrigger = time.Now()
	err = u.Update(trigger)
	if err != nil {
		return err
	}
	err = u.Emit(*trigger)
	if err != nil {
		return err
	}
	return nil
}

func (u *triggerService) EmitAll() error {
	all, err := u.repository.FindAll()
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

func (u *triggerService) FindById(id string) (*domain.Trigger, error) {
	return u.repository.FindById(id)
}

func (u *triggerService) Create(trigger *domain.Trigger) error {
	return u.repository.Create(trigger)
}

func (u *triggerService) Update(trigger *domain.Trigger) error {
	return u.repository.Update(trigger)
}

func (u *triggerService) Delete(trigger *domain.Trigger) error {
	return u.repository.Delete(trigger)
}
