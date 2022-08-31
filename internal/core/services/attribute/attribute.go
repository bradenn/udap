// Copyright (c) 2022 Braden Nicholson

package attribute

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/repository"
)

func New(db *gorm.DB) domain.AttributeService {
	repo := repository.NewAttributeRepository(db)
	operators := NewOperator()
	service := NewService(repo, operators)
	return service
}
