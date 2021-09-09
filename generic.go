package main

import (
	"github.com/google/uuid"
	"time"
)

type Persistent struct {
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time `sql:"index"`
	Id        uuid.UUID  `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
}
