// Copyright (c) 2022 Braden Nicholson

package generic

import "udap/internal/core/domain/common"

type Mock struct {
	common.Persistent
	Name  string `json:"name"`
	Value string `json:"value"`
}
