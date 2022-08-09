// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/core/domain"
)

type EntityRouter interface {
	RouteEntities(router chi.Router)
}

type entityRouter struct {
	service domain.EntityService
}

func NewEntityRouter(service domain.EntityService) EntityRouter {
	return &entityRouter{
		service: service,
	}
}

func (r entityRouter) RouteEntities(router chi.Router) {
	router.Route("/entities/{id}", func(local chi.Router) {
		local.Post("/icon", r.changeIcon)
	})
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
