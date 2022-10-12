// Copyright (c) 2022 Braden Nicholson

package domain

import (
	"fmt"
	"strings"
	"udap/internal/core/domain/common"
)

type ModuleConfig struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // Module, Daemon, etc.
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Variables   string `json:"variables"`
}

type Module struct {
	common.Persistent
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	UUID        string      `json:"uuid"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
	Author      string      `json:"author"`
	Variables   string      `json:"variables"`
	Channel     chan Module `json:"-" gorm:"-"`
	Config      string      `json:"config" gorm:"default:'{}'"`
	State       string      `json:"state"`
	Running     bool        `json:"running" gorm:"default:false"`
	Enabled     bool        `json:"enabled" gorm:"default:true"`
	Recover     int         `json:"recover"`
}

func (m *Module) SessionId() string {
	if m.UUID == "" {
		return "invalid"
	}
	return strings.Split(m.UUID, "-")[0]
}

func (m *Module) CompiledPath() string {
	if m.UUID == "" {
		return "invalid"
	}
	return strings.Replace(m.Path, ".go", fmt.Sprintf("-%s.so", m.UUID), 1)
}
