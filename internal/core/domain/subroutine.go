// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type SubRoutine struct {
	common.Persistent
	TriggerId   string  `json:"triggerId"`
	Macros      []Macro `json:"macros" gorm:"many2many:subroutine_macros;"`
	Description string  `json:"description"`
}
