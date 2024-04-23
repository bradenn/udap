// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type TriggerRepository interface {
	common.Persist[domain.Trigger]
	FindByName(name string) (*domain.Trigger, error)
	Register(*domain.Trigger) error
}

type TriggerOperator interface {
	Run(trigger domain.Trigger) error
	RunCustom(trigger domain.Trigger, key string, value string) error
}

type TriggerService interface {
	domain.Observable
	Trigger(name string) error
	TriggerCustom(name string, key string, value string) error

	Register(*domain.Trigger) error
	FindById(id string) (*domain.Trigger, error)
	Create(*domain.Trigger) error
	Update(*domain.Trigger) error
	Delete(*domain.Trigger) error
}
