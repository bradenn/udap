// Copyright (c) 2022 Braden Nicholson

package logs

import (
	"udap/internal/core/ports"
	"udap/internal/core/repository"
)

func New() ports.LogService {
	repo := repository.NewLogsRepository()
	service := NewService(repo)
	return service
}
