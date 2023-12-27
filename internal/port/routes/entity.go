// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"udap/internal/core/domain"
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
		local.Post("/update", r.update)
	})
}

func (r entityRouter) RouteExternal(_ chi.Router) {
}

func (r entityRouter) update(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(w, "invalid icon body", 401)
		return
	}
	e := domain.Entity{}
	err = json.Unmarshal(buf.Bytes(), &e)
	if err != nil {
		return
	}
	id := chi.URLParam(req, "id")
	if id != "" {
		err = r.service.ChangeAlias(id, e.Alias)
		if err != nil {
			http.Error(w, "invalid entity id", 401)
			return
		}
		err = r.service.ChangeIcon(id, e.Icon)
		if err != nil {
			http.Error(w, "invalid entity id", 401)
			return
		}
	}
	w.WriteHeader(200)
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
