// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type SubRoutineRepository interface {
	common.Persist[domain.SubRoutine]
	FindByTriggerId(id string) (res []*domain.SubRoutine, err error)
	AddMacro(subroutine *domain.SubRoutine, id string) error
	RemoveMacro(subroutine *domain.SubRoutine, id string) error
}

type SubRoutineOperator interface {
	Run(routine domain.SubRoutine) error
}

type SubRoutineService interface {
	domain.Observable
	Run(id string) error
	TriggerById(id string) error
	FindById(id string) (*domain.SubRoutine, error)
	AddMacro(id string, macroId string) error
	RemoveMacro(id string, macroId string) error
	Create(*domain.SubRoutine) error
	Update(*domain.SubRoutine) error
	Delete(id string) error
}
