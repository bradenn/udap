// Copyright (c) 2022 Braden Nicholson

package logs

import (
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

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

func NewService(repository ports.LogRepository) ports.LogService {
	return &logService{repository: repository}
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
