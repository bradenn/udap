// Copyright (c) 2021 Braden Nicholson

package models

import (
	"udap/internal/store"
	"udap/pkg/plugin"
)

type Daemon struct {
	store.Persistent
	plugin.Plugin `gorm:"-"`
	plugin.Metadata
	// Path refers to the literal name of the module
	Path string `json:"path"`
}
