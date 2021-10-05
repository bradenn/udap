package main

import (
	"github.com/google/uuid"
	"time"
)

type Persistent struct {
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time `sql:"index"`
	// Id is primary key of the persistent type, represented as a UUIDv4
	Id uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
}
