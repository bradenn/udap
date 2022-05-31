// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type User struct {
	common.Persistent
	Username string `json:"username"`
	First    string `json:"first"`
	Middle   string `json:"middle"`
	Last     string `json:"last"`
	Type     string `json:"type"`
	Password string `json:"password"`
}

type UserRepository interface {
	common.Persist[User]
}

type UserService interface {
	Observable
	Register(*User) error
	Authenticate(*User) error
	FindAll() (*[]User, error)
	FindById(id string) (*User, error)
	Create(*User) error
	FindOrCreate(*User) error
	Update(*User) error
	Delete(*User) error
}
