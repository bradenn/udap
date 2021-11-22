// Copyright (c) 2021 Braden Nicholson

package udap

import (
	"time"
)

type Persistent struct {
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time `sql:"index"`
	// Id is primary key of the persistent type, represented as a UUIDv4
	Id string `json:"id" gorm:"primary_key;type:string;default:uuid_generate_v4()"`
}

func (e *Persistent) UUID() string {
	return e.Id
}
