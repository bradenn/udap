// Copyright (c) 2022 Braden Nicholson

package macro

import (
	"gorm.io/gorm"
	"udap/internal/controller"
	"udap/internal/core/operators"
	"udap/internal/core/ports"
	"udap/internal/core/repository"
)

func New(db *gorm.DB, controller *controller.Controller) ports.MacroService {
	repo := repository.NewMacroRepository(db)
	operator := operators.NewMacroOperator(controller)
	service := NewService(repo, operator)

	return service
}
