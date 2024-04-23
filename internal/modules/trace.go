// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/port/routes"
	"udap/internal/srv"
)

func NewTrace(sys srv.System) {
	sys.WithRoute(routes.NewTraceRouter(sys.Store()))
}
