// Copyright (c) 2022 Braden Nicholson

package routes

import "github.com/go-chi/chi"

type Routable interface {
	RouteInternal(router chi.Router)
	RouteExternal(router chi.Router)
}
