// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"sync"
	"udap/internal/bond"
	"udap/internal/log"
	"udap/internal/models"
	"udap/internal/store"
)

type Users struct {
	PolyBuffer
	Observable
}

// Handle routes websocket requests to the appropriate function
func (u *Users) Handle(event bond.Msg) (res interface{}, err error) {
	switch event.Operation {
	case "register":
		return u.register(event)
	case "authenticate":
		return u.authenticate(event)
	}
	return nil, nil
}

func (u *Users) EmitAll() (err error) {

	for _, k := range u.Keys() {
		find := u.Find(k)
		u.emit(k, find)
	}

	return nil
}

func (u *Users) FetchAll() []models.User {
	var users []models.User
	log.Log("Fetching")
	err := store.DB.Table("users").Find(&users).Error
	if err != nil {
		return nil
	}
	return users
}

func LoadUsers() (m *Users) {
	m = &Users{}
	m.data = sync.Map{}
	m.raw = map[string]interface{}{}
	m.Run()
	m.FetchAll()
	return m
}

// register will create a new user from the provided body within 'payload'
func (u *Users) register(msg bond.Msg) (res interface{}, err error) {
	user := models.User{}
	err = user.Parse([]byte(msg.Payload))
	if err != nil {
		return nil, err
	}
	err = user.Register(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// authenticate will attempt to verify a user password and username combination
func (u *Users) authenticate(msg bond.Msg) (res interface{}, err error) {
	user := models.User{}
	err = user.Parse([]byte(msg.Payload))
	if err != nil {
		return nil, err
	}
	dbUser, err := user.Authenticate(user)
	if err != nil {
		return nil, err
	}
	return dbUser, nil
}

// Pull is the level at which this service needs to run
func (u *Users) Pull() {
	for _, k := range u.Keys() {
		err := u.get(k)
		if err != nil {
			return
		}
	}
}

func (u *Users) Find(name string) *models.User {
	return u.get(name).(*models.User)
}

func (u *Users) Set(id string, entity *models.User) {
	u.set(id, entity)
	u.emit(id, entity)
}
