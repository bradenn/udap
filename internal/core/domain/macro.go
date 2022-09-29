// Copyright (c) 2022 Braden Nicholson

package domain

import (
	"udap/internal/core/domain/common"
)

type Macro struct {
	common.Persistent
	Name        string `json:"name"`
	Description string `json:"description"`
	ZoneId      string `json:"zone"`
	Type        string `json:"type"`
	Value       string `json:"value"`
}

type MacroRepository interface {
	common.Persist[Macro]
}

type MacroOperator interface {
	Run(macro Macro) error
}

type MacroService interface {
	Observable
	FindAll() (*[]Macro, error)
	Run(id string) error
	FindById(id string) (*Macro, error)
	Create(*Macro) error
	FindOrCreate(*Macro) error
	Update(*Macro) error
	Delete(*Macro) error
}
