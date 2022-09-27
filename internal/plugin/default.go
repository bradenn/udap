// Copyright (c) 2021 Braden Nicholson

package plugin

import (
	"fmt"
	"time"
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/log"
)

type Config struct {
	Name        string     `json:"name"`
	Type        string     `json:"type"` // Module, Daemon, etc.
	Description string     `json:"description"`
	Version     string     `json:"version"`
	Author      string     `json:"author"`
	Variables   []Variable `json:"variables"`
}

type Variable struct {
	Name        string `json:"name"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

type Module struct {
	Config
	LastUpdate time.Time
	Frequency  time.Duration
	UUID       string
	*controller.Controller
}

// LogF is called once at the launch of the module
func (m *Module) LogF(format string, args ...any) {
	out := domain.Log{
		Group:   "module",
		Level:   "info",
		Event:   m.Name,
		Time:    time.Now(),
		Message: fmt.Sprintf(format, args...),
	}
	log.Event("%s::%s %s", out.Group, out.Event, out.Message)
	err := m.Logs.Create(&out)
	if err != nil {
		return
	}
}

func (m *Module) WarnF(format string, args ...any) {
	out := domain.Log{
		Group:   "module",
		Level:   "warn",
		Event:   m.Name,
		Time:    time.Now(),
		Message: fmt.Sprintf(format, args...),
	}
	log.Event("%s::%s %s", out.Group, out.Event, out.Message)
	err := m.Logs.Create(&out)
	if err != nil {
		return
	}
}

func (m *Module) ErrF(format string, args ...any) {
	out := domain.Log{
		Group:   "module",
		Level:   "error",
		Event:   m.Name,
		Time:    time.Now(),
		Message: fmt.Sprintf(format, args...),
	}
	log.Event("%s::%s %s", out.Group, out.Event, out.Message)
	err := m.Logs.Create(&out)
	if err != nil {
		return
	}
}

func (m *Module) Err(err error) {
	if err == nil {
		return
	}
	out := domain.Log{
		Group:   "module",
		Level:   "error",
		Event:   m.Name,
		Time:    time.Now(),
		Message: fmt.Sprintf("Error: %s", err.Error()),
	}
	log.Event("%s::%s %s", out.Group, out.Event, out.Message)
	err = nil
	err = m.Logs.Create(&out)
	if err != nil {
		return
	}
}

// UpdateInterval is called once at the launch of the module
func (m *Module) UpdateInterval(frequency time.Duration) error {
	m.Frequency = time.Millisecond * frequency
	m.LastUpdate = time.Now().Add(-m.Frequency)
	return nil
}

// Ready is used to determine whether the module should update
func (m *Module) Ready() bool {
	if time.Since(m.LastUpdate) >= m.Frequency {
		m.LastUpdate = time.Now()
		return true
	}
	return false
}

// Connect is called once at the launch of the module
func (m *Module) Connect(ctrl *controller.Controller, uuid string) error {
	m.LastUpdate = time.Now()
	m.Controller = ctrl
	m.UUID = uuid
	return nil
}

func (m *Module) OnEmit() error {
	return nil
}

func (m *Module) InitConfig(key string, value string) error {
	return m.Modules.InitConfig(m.UUID, key, value)
}

func (m *Module) GetConfig(key string) (string, error) {
	return m.Modules.GetConfig(m.UUID, key)
}

func (m *Module) SetConfig(key string, value string) error {
	return m.Modules.SetConfig(m.UUID, key, value)
}

// Dispose is called once at the launch of the module
func (m *Module) Dispose() error {

	return nil
}
