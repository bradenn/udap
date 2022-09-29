// Copyright (c) 2022 Braden Nicholson

package macro

import (
	"udap/internal/core/domain"
	"udap/internal/core/generic"
)

type macroService struct {
	repository domain.MacroRepository
	operator   domain.MacroOperator
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

func NewService(repository domain.MacroRepository, operator domain.MacroOperator) domain.MacroService {
	return &macroService{repository: repository, operator: operator}
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

func (u *macroService) FindOrCreate(macro *domain.Macro) error {
	return u.repository.FindOrCreate(macro)
}

func (u *macroService) Update(macro *domain.Macro) error {
	return u.repository.Update(macro)
}

func (u *macroService) Delete(macro *domain.Macro) error {
	return u.repository.Delete(macro)
}
