// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type User struct {
	common.Persistent
	Username string `json:"username"`
	First    string `json:"first"`
	Middle   string `json:"middle"`
	Last     string `json:"last"`
	Type     string `json:"type"`
	Password string `json:"password"`
}
