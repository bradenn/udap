// Copyright (c) 2022 Braden Nicholson

package module

import (
	"gorm.io/gorm"
	"udap/internal/controller"
	"udap/internal/core/operators"
	"udap/internal/core/ports"
	"udap/internal/core/repository"
)

func New(db *gorm.DB, controller *controller.Controller) ports.ModuleService {
	repo := repository.NewModuleRepository(db)
	operator := operators.NewModuleOperator(controller)
	service := NewService(repo, operator)
	return service
}
