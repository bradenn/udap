// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type UserRepository interface {
	common.Persist[domain.User]
}

type UserService interface {
	domain.Observable
	Register(*domain.User) error
	Authenticate(*domain.User) error
	FindAll() (*[]domain.User, error)
	FindById(id string) (*domain.User, error)
	Create(*domain.User) error
	FindOrCreate(*domain.User) error
	Update(*domain.User) error
	Delete(*domain.User) error
}
