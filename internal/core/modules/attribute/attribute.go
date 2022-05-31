// Copyright (c) 2022 Braden Nicholson

package attribute

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

func New(db *gorm.DB) domain.AttributeService {
	repo := NewRepository(db)
	operators := NewOperator()
	service := NewService(repo, operators)
	return service
}
