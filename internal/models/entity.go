// Copyright (c) 2021 Braden Nicholson

package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
	"udap/internal/cache"
	"udap/internal/store"
)

type iEntity interface {
	GetAttributes() ([]Attribute, error)
	SetAttribute(key string, value Attribute) error
	GetAttribute(key string) (Attribute, error)
}

type Entity struct {
	store.Persistent
	LastPoll   time.Time    `json:"lastPoll"`
	Name       string       `gorm:"unique" json:"name"`               // Given name from module
	Alias      string       `json:"alias"`                            // Name from users
	Type       string       `json:"type"`                             // Type of entity {Light, Sensor, Etc}
	Module     string       `json:"module"`                           // Parent Module name
	Neural     string       `json:"neural" gorm:"default:'inactive'"` // Parent Module name
	Locked     bool         `json:"locked"`                           // Is the Entity state locked?
	Protocol   string       `json:"protocol"`                         // scalar
	Attributes []*Attribute `json:"attributes" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:entity_id"`
	Icon       string       `json:"icon" gorm:"default:'􀛮'"` // The icon to represent this entity
	Frequency  int          `json:"frequency" gorm:"default:3000"`
	Predicted  string       `gorm:"-" json:"predicted"` // scalar
	State      string       `json:"state"`
	Config     string       `json:"config"`
	Live       bool         `gorm:"-" json:"live"`
}

func (e *Entity) BeforeCreate(tx *gorm.DB) (err error) {

	return
}

func (e *Entity) AfterFind(tx *gorm.DB) (err error) {

	return
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

func (e *Entity) OnChange(fn func(entity Entity) error) error {
	cache.WatchFn(fmt.Sprintf("entity.%s", e.Id), func(s string) error {
		p := Entity{}
		err := json.Unmarshal([]byte(s), &p)
		if err != nil {
			return err
		}
		err = fn(p)
		if err != nil {
			return err
		}
		return nil
	})
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

func (e *Entity) pull() error {
	for _, attribute := range e.Attributes {
		if attribute.Key == "on" {
			err := attribute.Get()
			if err != nil {
				return err
			}
		}
	}
	return nil
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
	err = e.writeCache()
	if err != nil {
		return err
	}

	return err
}

func (e *Entity) writeCache() error {
	marshal, err := json.Marshal(e)
	if err != nil {
		return err
	}
	if marshal == nil {
		return nil
	}
	err = cache.PutLn(string(marshal), "entity", e.Id)
	if err != nil {
		return err
	}
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Entity) readCache() error {
	ln, err := cache.GetLn("entity", e.Id)
	if err != nil {
		return err
	}
	s := ln.(string)
	if s == "" {
		err = e.writeCache()
		if err != nil {
			return err
		}
		return nil
	}
	err = json.Unmarshal([]byte(s), &e)
	if err != nil {
		return err
	}
	return nil
}

type Tag struct {
	Attribute string          `json:"attribute"`
	Value     json.RawMessage `json:"value"`
}

func (e *Entity) Push(state string) error {
	t := Tag{}
	err := json.Unmarshal([]byte(state), &t)
	if err != nil {
		return err
	}

	for _, attribute := range e.Attributes {
		if strings.ToLower(attribute.Key) == strings.ToLower(t.Attribute) {
			err = attribute.Put(t.Value)
			if err != nil {
				return err
			}
			err = e.writeCache()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Entity) Logs() ([]Log, error) {
	var lgs []Log
	err := store.DB.Model(&Log{}).Where("entity_id = ? AND cct <> 0", e.Id).Find(&lgs).Error
	if err != nil {
		return nil, err
	}
	return lgs, nil
}

func (e *Entity) timestamp() {
	e.LastPoll = time.Now()
}

func (e *Entity) coolDown() bool {
	return time.Since(e.LastPoll) < time.Millisecond*time.Duration(e.Frequency)
}

func (e *Entity) Pull() error {
	if e.coolDown() {
		return nil
	}

	e.timestamp()
	err := e.pull()
	if err != nil {
		return err
	}

	return nil
}
