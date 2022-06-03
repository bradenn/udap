// Copyright (c) 2022 Braden Nicholson

package common

import "time"

type Persistent struct {
	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"updated"`
	Deleted   bool       `json:"deleted"`
	deletedAt *time.Time `sql:"index"`
	Id        string     `json:"id" gorm:"primary_key;type:string;default:uuid_generate_v4()"`
}

type Persist[T any] interface {
	FindAll() (*[]T, error)
	FindById(id string) (*T, error)
	Create(*T) error
	FindOrCreate(*T) error
	Update(*T) error
	Delete(*T) error
}
