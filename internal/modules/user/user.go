// Copyright (c) 2022 Braden Nicholson

package user

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

func New(db *gorm.DB) domain.UserService {
	repo := NewUserRepository(db)
	service := NewUserService(repo)
	return service
}
