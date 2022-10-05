// Copyright (c) 2022 Braden Nicholson

package operators

import (
	"udap/internal/core/domain"
	"udap/internal/core/ports"
)

type notificationOperator struct {
}

func NewNotificationOperator() ports.NotificationOperator {
	return &notificationOperator{}
}

func (m *notificationOperator) Send(notification domain.Notification) error {

	return nil
}
