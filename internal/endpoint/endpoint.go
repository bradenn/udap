// Copyright (c) 2022 Braden Nicholson

package endpoint

import (
	"gorm.io/gorm"
	"math/rand"
	"time"
	"udap/internal/store"
)

type Endpoint struct {
	store.Persistent
	Name string `json:"name" gorm:"unique"`
	Key  string `json:"key"`
}

// BeforeCreate is a hook function from gorm, called when an endpoint is inserted
func (e *Endpoint) BeforeCreate(_ *gorm.DB) error {
	e.Key = randomSequence()
	return nil
}

// findById attempts to find an endpoint by its UUID
func findById(id string) (endpoint Endpoint, err error) {
	db := store.DB.Where("id = ?", id).First(&endpoint)
	err = db.Error
	return
}

// insert will attempt to create a database record if one does not already exist
func insert(endpoint *Endpoint) (err error) {
	db := store.DB.Create(endpoint)
	err = db.Error
	return
}

// insert will attempt to create a database record if one does not already exist
func drop(endpoint *Endpoint) (err error) {
	db := store.DB.Delete(endpoint)
	err = db.Error
	return
}

// randomSequence generates a random id for use as an endpoint key
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
