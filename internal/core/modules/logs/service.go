// Copyright (c) 2022 Braden Nicholson

package logs

import (
	"fmt"
	"udap/internal/core/domain"
)

type logService struct {
	repository domain.LogRepository
	channel    chan<- domain.Mutation
}

func (u *logService) EmitAll() error {
	all, err := u.repository.FindAll()
	if err != nil {
		return err
	}
	for _, log := range all {
		err = u.emit(&log)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *logService) emit(log *domain.Log) error {
	if u.channel == nil {
		return nil
	}
	u.channel <- domain.Mutation{
		Status:    "update",
		Operation: "log",
		Body:      *log,
		Id:        log.Id,
	}
	return nil
}

func (u *logService) Watch(mut chan<- domain.Mutation) error {
	if u.channel != nil {
		return fmt.Errorf("channel already set")
	}
	u.channel = mut

	return nil
}

func NewService(repository domain.LogRepository) domain.LogService {
	return &logService{repository: repository}
}

// Repository Mapping

func (u *logService) Create(log *domain.Log) error {
	err := u.repository.Create(log)
	if err != nil {
		return err
	}
	err = u.emit(log)
	if err != nil {
		return err
	}
	return nil
}
