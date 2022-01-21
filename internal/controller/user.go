// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"udap/internal/bond"
	"udap/internal/models"
)

type UserController struct {
	users models.Users
}

// Handle routes websocket requests to the appropriate function
func (d *UserController) Handle(event bond.Msg) (res any, err error) {
	switch event.Operation {
	case "register":
		return d.register(event)
	case "authenticate":
		return d.authenticate(event)
	}
	return nil, nil
}

// register will create a new user from the provided body within 'payload'
func (d *UserController) register(msg bond.Msg) (res any, err error) {
	user := models.User{}
	err = user.Parse([]byte(msg.Payload))
	if err != nil {
		return nil, err
	}
	err = d.users.Register(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// authenticate will attempt to verify a user password and username combination
func (d *UserController) authenticate(msg bond.Msg) (res any, err error) {
	user := models.User{}
	err = user.Parse([]byte(msg.Payload))
	if err != nil {
		return nil, err
	}
	dbUser, err := d.users.Authenticate(user)
	if err != nil {
		return nil, err
	}
	return dbUser, nil
}
