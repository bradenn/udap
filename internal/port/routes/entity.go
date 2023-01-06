// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/core/ports"
)

type entityRouter struct {
	service ports.EntityService
}

func NewEntityRouter(service ports.EntityService) Routable {
	return &entityRouter{
		service: service,
	}
}

func (r entityRouter) RouteInternal(router chi.Router) {
	router.Route("/entities/{id}", func(local chi.Router) {
		local.Post("/icon", r.changeIcon)
		local.Post("/alias", r.changeAlias)
	})
}

func (r entityRouter) RouteExternal(_ chi.Router) {
}

func (r entityRouter) changeAlias(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(w, "invalid icon body", 401)
		return
	}
	alias := buf.String()
	id := chi.URLParam(req, "id")
	if id != "" {
		err = r.service.ChangeAlias(id, alias)
		if err != nil {
			http.Error(w, "invalid entity id", 401)
			return
		}
	}
	w.WriteHeader(200)
}

func (r entityRouter) changeIcon(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(w, "invalid icon body", 401)
		return
	}
	icon := buf.String()
	id := chi.URLParam(req, "id")
	if id != "" {
		err = r.service.ChangeIcon(id, icon)
		if err != nil {
			http.Error(w, "invalid entity id", 401)
			return
		}
	}
	w.WriteHeader(200)
}
