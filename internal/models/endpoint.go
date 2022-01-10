// Copyright (c) 2021 Braden Nicholson

package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"time"
	"udap/internal/log"
	"udap/internal/store"
)

// var m *sync.RWMutex
//
// func init() {
// 	m = &sync.RWMutex{}
// }

// Endpoint represents a client device connected to the UDAP network
type Endpoint struct {
	store.Persistent

	Name string `json:"name" gorm:"unique"`

	Type string `json:"type"`

	Remote *Remote `gorm:"-"`

	Frequency int `json:"frequency" gorm:"default:3000"`

	key string

	registered bool

	enrolledSince time.Time `gorm:"-"`
}

func (e *Endpoint) Compile() (a map[string]any) {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(marshal, &a)
	if err != nil {
		return nil
	}
	return a
}

func (e *Endpoint) Enroll() error {
	err := store.DB.Model(&Endpoint{}).FirstOrCreate(&e).Error
	if err != nil {
		return err
	}

	e.registered = true
	e.enrolledSince = time.Now()
	log.Log("%s: Enrolled (%s)", e.Name)

	return nil
}

func (e *Endpoint) Enrolled() bool {
	return e.registered
}

// BeforeCreate is a hook function from gorm, called when an endpoint is inserted
func (e *Endpoint) BeforeCreate(_ *gorm.DB) error {
	e.key = randomSequence()
	return nil
}

// randomSequence generates a random id for use as a key
func randomSequence() string {
	template := "abcdefghijklmnopqrstuvwxyz"
	var out string
	rand.Seed(time.Now().Unix())
	for i := 0; i < 8; i++ {
		r := rand.Intn(26)
		u := template[r]
		out += string(u)
	}
	return out
}

// Fetch populates the struct from the database
func (e *Endpoint) Fetch() error {
	if e.Id == "" {
		return fmt.Errorf("invalid id")
	}
	err := store.DB.Where("name = ? OR id = ?", e.Name, e.Id).First(&e).Error
	if err != nil {
		return err
	}

	return nil
}

func (e *Endpoint) Unenroll() {
	e.registered = false
	log.Log("%s: Unenrolled... (%s)", e.Name, time.Since(e.enrolledSince).String())
}
