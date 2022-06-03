// Copyright (c) 2022 Braden Nicholson

package user

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
)

type userRepo struct {
	generic.Store[domain.User]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.UserRepository {
	return &userRepo{
		db:    db,
		Store: generic.NewStore[domain.User](db),
	}
}
