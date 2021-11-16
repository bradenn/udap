package types

import (
	"github.com/google/uuid"
)

type Switch struct {
	Persistent
	InstanceId uuid.UUID `json:"instance_id"`
	Name       string    `json:"name"`
	Zone       string    `json:"zone"`
	Type       string    `json:"type"`
}
