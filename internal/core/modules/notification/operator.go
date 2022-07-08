// Copyright (c) 2022 Braden Nicholson

package notification

import (
	"udap/internal/controller"
	"udap/internal/core/domain"
)

type notificationOperator struct {
	ctrl *controller.Controller
}

func NewOperator(ctrl *controller.Controller) domain.NotificationOperator {
	return &notificationOperator{
		ctrl: ctrl,
	}
}

func (m *notificationOperator) Send(notification domain.Notification) error {
	err := m.ctrl.Notifications.Create(&notification)
	if err != nil {
		return err
	}
	return nil
}
