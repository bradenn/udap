// Copyright (c) 2022 Braden Nicholson

package domain

import (
	"time"
	"udap/internal/core/domain/common"
)

const (
	SYSTEM = "system"
	MODULE = "module"
	MANUAL = "manual"
)

type Trigger struct {
	common.Persistent
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	LastTrigger time.Time `json:"lastTrigger"`
}

func NewTrigger(name string, description string, triggerType string) Trigger {
	return Trigger{
		Name:        name,
		Type:        triggerType,
		Description: description,
	}
}
