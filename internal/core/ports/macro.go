// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"time"
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type MacroRepository interface {
	common.Persist[domain.Macro]
}

type MacroOperator interface {
	Run(macro domain.Macro) error
	RunAndRevert(macro domain.Macro, baseline domain.Macro, revert time.Duration) error
}

type MacroService interface {
	domain.Observable
	FindAll() (*[]domain.Macro, error)
	Run(id string) error
	RunAndRevert(id string, revert time.Duration) error
	FindById(id string) (*domain.Macro, error)
	Create(*domain.Macro) error
	Update(*domain.Macro) error
	Delete(id string) error
}
