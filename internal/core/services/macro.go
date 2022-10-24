// Copyright (c) 2022 Braden Nicholson

package services

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
	"udap/internal/core/repository"
)

func NewMacroService(db *gorm.DB, operator ports.MacroOperator) ports.MacroService {
	repo := repository.NewMacroRepository(db)

	return &macroService{repository: repo, operator: operator}
}

type macroService struct {
	repository ports.MacroRepository
	operator   ports.MacroOperator
	generic.Watchable[domain.Macro]
}

func (u *macroService) Run(id string) error {
	byId, err := u.FindById(id)
	if err != nil {
		return err
	}
	err = u.operator.Run(*byId)
	if err != nil {
		return err
	}
	return nil
}

func (u *macroService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	for _, macro := range *all {
		err = u.Emit(macro)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *macroService) mutate(macro *domain.Macro) error {
	err := u.repository.Update(macro)
	if err != nil {
		return err
	}
	err = u.Emit(*macro)
	if err != nil {
		return err
	}
	return nil
}

// Repository Mapping

func (u *macroService) FindAll() (*[]domain.Macro, error) {
	return u.repository.FindAll()
}

func (u *macroService) FindById(id string) (*domain.Macro, error) {
	return u.repository.FindById(id)
}

func (u *macroService) Create(macro *domain.Macro) error {
	err := u.repository.Create(macro)
	if err != nil {
		return err
	}
	err = u.Emit(*macro)
	if err != nil {
		return err
	}
	return nil
}

func (u *macroService) Update(macro *domain.Macro) error {
	err := u.repository.Update(macro)
	if err != nil {
		return err
	}
	err = u.Emit(*macro)
	if err != nil {
		return err
	}
	return nil
}

func (u *macroService) Delete(id string) error {
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
