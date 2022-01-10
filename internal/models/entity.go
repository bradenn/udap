// Copyright (c) 2021 Braden Nicholson

package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"udap/internal/cache"
	"udap/internal/store"
)

type LightState struct {
	Power string `json:"power"`
	Red   int    `json:"red"`
	Green int    `json:"green"`
	Blue  int    `json:"blue"`
	Level int    `json:"level"`
	CCT   int    `json:"cct"`
	Mode  string `json:"mode"`
}

func (l *LightState) Parse(msg json.RawMessage) {
	err := json.Unmarshal(msg, l)
	if err != nil {
		return
	}
}

func (l *LightState) JSON() json.RawMessage {
	bytes, err := json.Marshal(l)
	if err != nil {
		return nil
	}
	return bytes
}

func (l *LightState) IsOn() bool {
	return l.Power == "on"
}

type State json.RawMessage

type Entity struct {
	store.Persistent
	LastPoll  time.Time       `json:"lastPoll"`
	Name      string          `gorm:"unique" json:"name"`      // Given name from module
	Alias     string          `json:"alias"`                   // Name from users
	Type      string          `json:"type"`                    // Type of entity {Light, Sensor, Etc}
	Module    string          `json:"module"`                  // Parent Module name
	Locked    bool            `json:"locked"`                  // Is the Entity state locked?
	Protocol  string          `json:"protocol"`                // scalar
	Icon      string          `json:"icon" gorm:"default:'􀛮'"` // The icon to represent this entity
	Frequency int             `json:"frequency" gorm:"default:3000"`
	Predicted string          `gorm:"-" json:"predicted"` // scalar
	State     json.RawMessage `gorm:"-" json:"state"`
	Live      bool            `gorm:"-" json:"live"`
	tx        Tx
	rx        Rx
}

func (e *Entity) Unlock() error {
	if !e.Locked {
		return fmt.Errorf("this entity is not locked")
	}
	e.Locked = false
	err := e.Update()
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
	err := e.Update()
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) ChangeIcon(icon string) error {
	e.Icon = icon
	err := e.Update()
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
	err := e.Update()
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) Suggest(state string) error {
	e.Predicted = state
	err := cache.PutLn(state, e.Id, "suggested")
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

type Tx func(state State) error
type Rx func() State

// Handlers attempts to locate
func (e *Entity) Handlers(tx Tx, rx Rx) error {
	e.rx = rx
	e.Live = true
	err := e.pull()
	if err != nil {
		return err
	}
	e.tx = tx

	return nil
}

func (e *Entity) Connected() bool {
	return e.rx != nil && e.tx != nil
}

func (e *Entity) push() error {
	if e.tx == nil {
		return fmt.Errorf("entity '%s' is not connected to its parent module '%s'", e.Name, e.Module)
	}
	if e.Locked {
		return fmt.Errorf("entity '%s' is locked. Unlock it before making changes", e.Name)
	}
	err := e.tx(State(e.State))
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) pull() error {

	if e.rx == nil {
		return fmt.Errorf("rx is not set")
	}
	e.State = json.RawMessage(e.rx())
	return nil
}

func (e *Entity) Emplace() error {
	if e.Id == "" {
		err := store.DB.Model(&Entity{}).Where("name = ? AND module = ?", e.Name, e.Module).FirstOrCreate(&e).Error
		if err != nil {
			return err
		}
	} else {
		err := store.DB.Model(&Entity{}).Where("id = ?", e.Name).First(&e).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Entity) Delete() error {
	err := store.DB.Where("name = ? AND module = ?", e.Name, e.Module).Delete(&e).Error
	return err
}

func (e *Entity) Update() error {
	err := store.DB.Where("id = ?", e.Id).Save(&e).Error
	return err
}

func (e *Entity) writeCache() error {
	err := cache.PutLn(string(e.State), e.Id, "state")
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) readCache() error {
	ln, err := cache.GetLn(e.Id, "state")
	if err != nil {
		return err
	}

	e.State = json.RawMessage(ln.(string))
	return nil
}

func (e *Entity) Push(state State) error {

	e.State = json.RawMessage(state)
	e.timestamp()
	err := e.push()
	if err != nil {
		err = e.readCache()
		if err != nil {
			return err
		}
		return err
	}
	err = e.writeCache()
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) timestamp() {
	e.LastPoll = time.Now()
}

func (e *Entity) Pull() error {
	if time.Since(e.LastPoll) >= time.Millisecond*time.Duration(e.Frequency) {
		e.timestamp()
		err := e.pull()
		if err != nil {
			return err
		}
		err = e.writeCache()
		if err != nil {
			return err
		}
		return nil
	}

	ln2, err := cache.GetLn(e.Id, "suggested")
	if err != nil {
		e.Predicted = string(e.State)
	} else {
		e.Predicted = ln2.(string)
	}

	err = e.readCache()
	if err != nil {
		return err
	}

	return nil
}
