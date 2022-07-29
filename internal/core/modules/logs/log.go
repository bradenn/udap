// Copyright (c) 2022 Braden Nicholson

package logs

import (
	"udap/internal/core/domain"
)

func New() domain.LogService {
	repo := NewRepository()
	service := NewService(repo)
	return service
}
