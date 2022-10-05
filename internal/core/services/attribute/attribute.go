// Copyright (c) 2022 Braden Nicholson

package attribute

import (
	"gorm.io/gorm"
	"udap/internal/core/ports"
	"udap/internal/core/repository"
)

func New(db *gorm.DB) ports.AttributeService {
	repo := repository.NewAttributeRepository(db)
	operators := NewOperator()
	service := NewService(repo, operators)
	return service
}
