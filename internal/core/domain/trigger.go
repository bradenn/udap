// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type Trigger struct {
	common.Persistent
	Name        string `json:"name"`
	Description string `json:"description"`
}
