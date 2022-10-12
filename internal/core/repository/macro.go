// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

type macroRepository struct {
	generic.Store[domain.Macro]
	db *gorm.DB
}

func NewMacroRepository(db *gorm.DB) ports.MacroRepository {
	return &macroRepository{
		db:    db,
		Store: generic.NewStore[domain.Macro](db),
	}
}
