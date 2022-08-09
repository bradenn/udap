// Copyright (c) 2022 Braden Nicholson

package endpoint

import (
	"gorm.io/gorm"
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/core/repository"
)

func New(db *gorm.DB, controller *controller.Controller) domain.EndpointService {
	repo := repository.NewEndpointRepository(db)
	operator := NewOperator(controller)
	service := NewService(repo, operator)
	return service
}
