// Copyright (c) 2022 Braden Nicholson

package endpoint

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

func New(db *gorm.DB) domain.EndpointService {
	repo := NewRepository(db)
	operator := NewOperator()
	service := NewService(repo, operator)
	return service
}
