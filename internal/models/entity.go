// Copyright (c) 2021 Braden Nicholson

package models

import (
	"fmt"
	"strings"
	"time"
	"udap/internal/store"
)

type Entity struct {
	store.Persistent
	Name   string `gorm:"unique" json:"name"` // Given name from module
	Alias  string `json:"alias"`              // Name from users
	Type   string `json:"type"`               // Type of entity {Light, Sensor, Etc}
	Module string `json:"module"`             // Parent Module name
	Locked bool   `json:"locked"`             // Is the Entity state locked?
	Config string `json:"config"`

	Position string `json:"position" gorm:"default:'{}'"`

	LastPoll  time.Time `json:"lastPoll"`
	Neural    string    `json:"neural" gorm:"default:'inactive'"` // Parent Module name
	Protocol  string    `json:"protocol"`                         // scalar
	Icon      string    `json:"icon" gorm:"default:'􀛮'"`          // The icon to represent this entity
	Frequency int       `json:"frequency" gorm:"default:3000"`

	Predicted string `gorm:"-" json:"predicted"` // scalar
	State     string `json:"state"`
	Live      bool   `gorm:"-" json:"live"`
}

func (e *Entity) Unlock() error {
	if !e.Locked {
		return fmt.Errorf("this entity is not locked")
	}
	e.Locked = false
	err := e.update()
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) Lock() error {
	if e.Locked {
		return fmt.Errorf("this entity is already locked")
	}
	e.Locked = true
	err := e.update()
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) ChangeConfig(value string) error {
	e.Config = value
	err := e.update()
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) ChangeNeural(value string) error {
	e.Neural = value
	err := e.update()
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) ChangeIcon(icon string) error {
	e.Icon = icon
	err := e.update()
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) Rename(name string) error {
	if e.Alias == name {
		return fmt.Errorf("alias has not been changed")
	}
	var cnt int64
	store.DB.Where("alias = ?", name).Count(&cnt)
	if cnt >= 1 {
		return fmt.Errorf("alias is already in use")
	}
	err := e.update()
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) Suggest(state string) error {
	e.Predicted = state
	err := e.update()
	if err != nil {
		return err
	}
	return nil
}

// Find attempts to locate
func (e *Entity) Find() error {
	err := store.DB.Where("name = ? AND module = ?", e.Name, e.Module).First(&e).Error
	return err
}

// Path attempts to locate
func (e *Entity) Path() string {
	return strings.ToLower(fmt.Sprintf("%s.%s", e.Module, e.Name))
}

func (e *Entity) Emplace() error {
	if e.Id == "" {
		err := store.DB.Model(&Entity{}).Where("name = ? AND module = ?", e.Name, e.Module).FirstOrCreate(e).Error
		if err != nil {
			return err
		}
	} else {
		err := store.DB.Model(&Entity{}).Where("id = ?", e.Name).First(e).Error
		if err != nil {
			return err
		}
	}
	err := store.DB.Model(&Entity{}).Where("id = ?", e.Id).Save(e).Error
	if err != nil {
		return err
	}

	return nil
}

func (e *Entity) delete() error {
	err := store.DB.Where("name = ? AND module = ?", e.Name, e.Module).Delete(&e).Error
	return err
}

func (e *Entity) update() error {
	err := store.DB.Where("id = ?", e.Id).Save(e).Error
	if err != nil {
		return err
	}

	return err
}

func NewMediaEntity(name string, module string) *Entity {
	e := Entity{
		Name:   name,
		Type:   "media",
		Module: module,
	}
	return &e
}

func NewSpectrum(name string, module string) *Entity {

	e := Entity{
		Name:   name,
		Type:   "spectrum",
		Module: module,
	}
	return &e
}

func NewDimmer(name string, module string) *Entity {

	e := Entity{
		Name:   name,
		Type:   "dimmer",
		Module: module,
	}
	return &e
}

func NewSwitch(name string, module string) *Entity {

	e := Entity{
		Name:   name,
		Type:   "switch",
		Module: module,
	}
	return &e
}
