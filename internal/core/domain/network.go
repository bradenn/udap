// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type Network struct {
	common.Persistent
	Name   string `json:"name"`
	Dns    string `json:"dns"`
	Router string `json:"index"`
	Lease  string `json:"lease"`
	Mask   string `json:"mask"`
	Range  string `json:"range"`
}
