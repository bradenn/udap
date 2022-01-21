// Copyright (c) 2022 Braden Nicholson

package models

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	"udap/internal/store"
)

type User struct {
	store.Persistent
	Username string `json:"username"`
	First    string `json:"first"`
	Middle   string `json:"middle"`
	Last     string `json:"last"`
	Password string `json:"password"`
}

func (u *User) Parse(data []byte) error {
	if !json.Valid(data) {
		return fmt.Errorf("failed to parse invalid json for type 'user'")
	}
	err := json.Unmarshal(data, u)
	if err != nil {
		return fmt.Errorf("failed to parse json for type 'user': %s", err.Error())
	}
	return nil
}

type Users struct {
	db store.Database
}

func (u *Users) Load(db store.Database) {
	u.db = db
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *Users) Register(user *User) error {
	password, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = password
	err = u.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *Users) FindById(id string) (user User, err error) {
	err = u.db.Where("id = ?", id).First(&user).Error
	return
}

func (u *Users) Authenticate(user User) (User, error) {
	pass := user.Password
	err := u.db.Where("username = ?", user.Username).First(&user).Error
	if err != nil {
		return User{}, nil
	}
	hash := CheckPasswordHash(pass, user.Password)
	if !hash {
		return User{}, fmt.Errorf("invalid password")
	}
	return user, nil
}
