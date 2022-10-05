// Copyright (c) 2022 Braden Nicholson

package domain

import (
	"time"
	"udap/internal/core/domain/common"
)

type Log struct {
	common.Persistent
	Group   string    `json:"group"`
	Event   string    `json:"event"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	Level   string    `json:"level"`
}
