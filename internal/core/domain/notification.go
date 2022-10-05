// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type Notification struct {
	common.Persistent
	Title    string `json:"title"`
	Target   string `json:"string"`
	Module   string `json:"module"`
	Body     string `json:"body"`
	Priority int    `json:"priority"`
}
