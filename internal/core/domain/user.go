// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type User struct {
	common.Persistent
	Username string `json:"username"`
	First    string `json:"first"`
	Middle   string `json:"middle"`
	Last     string `json:"last"`

	Residency      string `json:"residency"`      // resident, periodic, transient
	Classification string `json:"classification"` // human, feline, caine
	Type           string `json:"type"`
	Photo          string `json:"photo"`
	Password       string `json:"password"`
}
