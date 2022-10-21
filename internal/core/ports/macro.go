// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type MacroRepository interface {
	common.Persist[domain.Macro]
}

type MacroOperator interface {
	Run(macro domain.Macro) error
}

type MacroService interface {
	domain.Observable
	FindAll() (*[]domain.Macro, error)
	Run(id string) error
	FindById(id string) (*domain.Macro, error)
	Create(*domain.Macro) error
	Update(*domain.Macro) error
	Delete(id string) error
}
