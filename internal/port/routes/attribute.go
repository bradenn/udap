// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/core/domain"
	"udap/internal/log"
)

type AttributeRouter interface {
	RouteAttributes(router chi.Router)
}

type attributeRouter struct {
	service domain.AttributeService
}

func NewAttributeRouter(service domain.AttributeService) AttributeRouter {
	return &attributeRouter{
		service: service,
	}
}

func (r *attributeRouter) RouteAttributes(router chi.Router) {
	router.Route("/entities/{id}/attributes/{key}", func(local chi.Router) {
		local.Post("/request", r.request)
	})
}

func (r *attributeRouter) request(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	key := chi.URLParam(req, "key")
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		w.WriteHeader(400)
	}
	log.Event("Request '%s' = %s", key, buf.String())
	if id != "" && key != "" {
		err := r.service.Request(id, key, buf.String())
		if err != nil {
			w.WriteHeader(400)
		}
	}
	w.WriteHeader(200)
}
