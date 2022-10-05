// Copyright (c) 2022 Braden Nicholson

package domain

import (
	"udap/internal/core/domain/common"
)

type Macro struct {
	common.Persistent
	Name        string `json:"name"`
	Description string `json:"description"`
	ZoneId      string `json:"zone"`
	Type        string `json:"type"`
	Value       string `json:"value"`
}
