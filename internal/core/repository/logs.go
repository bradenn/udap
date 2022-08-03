// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"math/rand"
	"time"
	"udap/internal/core/domain"
)

type logsRepo struct {
	logs map[string][]domain.Log
}

func NewLogsRepository() domain.LogRepository {
	return &logsRepo{
		logs: map[string][]domain.Log{},
	}
}

func randomSequence() string {
	template := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	var out string
	rand.Seed(time.Now().Unix())
	for i := 0; i < 32; i++ {
		r := rand.Intn(62)
		u := template[r]
		out += string(u)
	}
	return out
}

func (m *logsRepo) Create(logEvent *domain.Log) error {
	logEvent.Id = randomSequence()
	m.logs[logEvent.Event] = append(m.logs[logEvent.Event], *logEvent)
	return nil
}

func (m *logsRepo) FindAll() (logs []domain.Log, err error) {
	for _, lgs := range m.logs {
		logs = append(logs, lgs...)
	}
	return logs, nil
}
