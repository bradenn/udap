// Copyright (c) 2022 Braden Nicholson

package services

import (
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
	"udap/internal/core/repository"
)

func NewLogService() ports.LogService {
	repo := repository.NewLogsRepository()
	return &logService{repository: repo}
}

type logService struct {
	repository ports.LogRepository
	generic.Watchable[domain.Log]
}

func (u *logService) EmitAll() error {
	all, err := u.repository.FindAll()
	if err != nil {
		return err
	}
	for _, log := range all {
		err = u.Emit(log)
		if err != nil {
			return err
		}
	}
	return nil
}

// Repository Mapping

func (u *logService) Create(log *domain.Log) error {
	err := u.repository.Create(log)
	if err != nil {
		return err
	}
	err = u.Emit(*log)
	if err != nil {
		return err
	}
	return nil
}
