// Copyright (c) 2022 Braden Nicholson

package services

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
	"udap/internal/core/repository"
)

func NewSubRoutineService(db *gorm.DB, operator ports.SubRoutineOperator) ports.SubRoutineService {
	repo := repository.NewSubRoutineRepository(db)
	return &subRoutineService{repository: repo, operator: operator}
}

type subRoutineService struct {
	repository ports.SubRoutineRepository
	operator   ports.SubRoutineOperator
	generic.Watchable[domain.SubRoutine]
}

func (u *subRoutineService) AddMacro(id string, macroId string) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	err = u.repository.AddMacro(byId, macroId)
	if err != nil {
		return err
	}
	return nil
}

func (u *subRoutineService) RemoveMacro(id string, macroId string) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	err = u.repository.RemoveMacro(byId, macroId)
	if err != nil {
		return err
	}
	return nil
}

func (u *subRoutineService) TriggerById(id string) error {
	routines, err := u.repository.FindByTriggerId(id)
	if err != nil {
		return err
	}

	for _, routine := range routines {
		err = u.operator.Run(*routine)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *subRoutineService) Run(id string) error {
	subroutine, err := u.FindById(id)
	if err != nil {
		return err
	}
	err = u.operator.Run(*subroutine)
	if err != nil {
		return err
	}
	return nil
}

func (u *subRoutineService) EmitAll() error {
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

func (u *subRoutineService) FindById(id string) (*domain.SubRoutine, error) {
	return u.repository.FindById(id)
}

func (u *subRoutineService) Create(subRoutine *domain.SubRoutine) error {
	return u.repository.Create(subRoutine)
}

func (u *subRoutineService) Update(subRoutine *domain.SubRoutine) error {
	return u.repository.Update(subRoutine)
}

func (u *subRoutineService) Delete(id string) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}

	err = u.repository.Delete(byId)
	if err != nil {
		return err
	}
	byId.Deleted = true
	err = u.Emit(*byId)
	if err != nil {
		return err
	}
	return err
}
