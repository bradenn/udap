// Copyright (c) 2022 Braden Nicholson

package domain

import "time"

type Persistent struct {
	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"updated"`
	Deleted   bool       `json:"deleted"`
	deletedAt *time.Time `sql:"index"`
	Id        string     `json:"id" gorm:"primary_key;type:string;default:uuid_generate_v4()"`
}
