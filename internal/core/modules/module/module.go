// Copyright (c) 2022 Braden Nicholson

package module

import (
	"gorm.io/gorm"
	"udap/internal/controller"
	"udap/internal/core/domain"
)

func New(db *gorm.DB, controller *controller.Controller) domain.ModuleService {
	repo := NewRepository(db)
	operator := NewOperator(controller)
	service := NewService(repo, operator)
	return service
}
