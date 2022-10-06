// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type TriggerRepository interface {
	common.Persist[domain.Trigger]
}

type TriggerOperator interface {
	Run(trigger domain.Trigger) error
}

type TriggerService interface {
	domain.Observable
	Trigger(id string) error
	FindById(id string) (*domain.Trigger, error)
	Create(*domain.Trigger) error
	Update(*domain.Trigger) error
	Delete(*domain.Trigger) error
}
