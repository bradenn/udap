// Copyright (c) 2022 Braden Nicholson

package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"udap/internal/core/domain"
)

type userService struct {
	repository domain.UserRepository
	channel    chan<- domain.Mutation
}

func (u *userService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	for _, user := range *all {
		err = u.emit(&user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *userService) emit(user *domain.User) error {
	if u.channel == nil {
		return nil
	}
	u.channel <- domain.Mutation{
		Status:    "update",
		Operation: "user",
		Body:      *user,
		Id:        user.Id,
	}
	return nil
}

func (u *userService) Watch(mut chan<- domain.Mutation) error {
	if u.channel != nil {
		return fmt.Errorf("channel already set")
	}
	u.channel = mut

	return nil
}

func NewService(repository domain.UserRepository) domain.UserService {
	return &userService{repository: repository}
}

// Services

func (u userService) Register(user *domain.User) error {
	password, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = password
	err = u.repository.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (u userService) Authenticate(user *domain.User) error {
	ref, err := u.repository.FindById(user.Id)
	if err != nil {
		return err
	}
	hash := CheckPasswordHash(user.Password, ref.Password)
	if !hash {
		return fmt.Errorf("invalid password")
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Repository Mapping

func (u userService) FindAll() (*[]domain.User, error) {
	return u.repository.FindAll()
}

func (u userService) FindById(id string) (*domain.User, error) {
	return u.repository.FindById(id)
}

func (u userService) Create(user *domain.User) error {
	return u.repository.Create(user)
}

func (u userService) FindOrCreate(user *domain.User) error {
	return u.repository.FindOrCreate(user)
}

func (u userService) Update(user *domain.User) error {
	return u.repository.Update(user)
}

func (u userService) Delete(user *domain.User) error {
	return u.repository.Delete(user)
}
