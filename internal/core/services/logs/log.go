// Copyright (c) 2022 Braden Nicholson

package logs

import (
	"udap/internal/core/domain"
	"udap/internal/core/repository"
)

func New() domain.LogService {
	repo := repository.NewLogsRepository()
	service := NewService(repo)
	return service
}
