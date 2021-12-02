// Copyright (c) 2021 Braden Nicholson

package controllers

import "github.com/go-chi/chi"

func Route(router chi.Router) {
	router.Route("/endpoint", EndpointRouter)
}
