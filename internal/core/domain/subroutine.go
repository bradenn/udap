package domain

import (
	"time"
	"udap/internal/core/domain/common"
)

type SubRoutine struct {
	common.Persistent
	TriggerId   string        `json:"triggerId"`
	Icon        string        `json:"icon" gorm:"default:'􁏀'"`
	Group       string        `json:"group"`
	Macros      []Macro       `json:"macros" gorm:"many2many:subroutine_macros;"`
	Description string        `json:"description"`
	RevertAfter time.Duration `json:"revertAfter"`
	LastRun     time.Time     `json:"lastRun"`
}
