// Copyright (c) 2022 Braden Nicholson

package user

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

func New(db *gorm.DB) domain.UserService {
	repo := NewRepository(db)
	service := NewService(repo)
	return service
}
