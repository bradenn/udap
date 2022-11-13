// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/repository"
	"udap/internal/core/services"
	"udap/internal/srv"
)

func NewNotifications(sys srv.System) {
	// Initialize service
	service := services.NewNotificationService(
		repository.NewNotificationRepository(sys.DB()))
	// Enroll routes
	sys.WithWatch(service)
	// sys.WithRoute(routes.NewNotificationRouter(service))
	sys.Ctrl().Notifications = service
}
