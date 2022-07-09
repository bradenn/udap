// Copyright (c) 2022 Braden Nicholson

package notification

import (
	"udap/internal/core/domain"
)

type notificationOperator struct {
}

func NewOperator() domain.NotificationOperator {
	return &notificationOperator{}
}

func (m *notificationOperator) Send(notification domain.Notification) error {

	return nil
}
