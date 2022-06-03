// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/core/domain"
)

type ModuleRouter interface {
	RouteModules(router chi.Router)
}

type moduleRouter struct {
	service domain.ModuleService
}

func NewModuleRouter(service domain.ModuleService) ModuleRouter {
	return &moduleRouter{
		service: service,
	}
}

func (r *moduleRouter) RouteModules(router chi.Router) {
	router.Route("/modules/{id}", func(local chi.Router) {
		local.Post("/build", r.build)
		local.Post("/disable", r.disable)
		local.Post("/enable", r.enable)
		local.Post("/halt", r.halt)
	})
}

func (r moduleRouter) build(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id != "" {
		mod, err := r.service.FindByName(id)
		if err != nil {
			return
		}
		err = r.service.Build(mod)
		if err != nil {
			http.Error(w, "invalid module name", 401)
		}
	}
	w.WriteHeader(200)
}

func (r moduleRouter) enable(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id != "" {
		err := r.service.Enable(id)
		if err != nil {
			http.Error(w, "invalid module name", 401)
		}
	}
	w.WriteHeader(200)
}

func (r moduleRouter) disable(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id != "" {
		err := r.service.Disable(id)
		if err != nil {
			http.Error(w, "invalid module name", 401)
		}
	}
	w.WriteHeader(200)
}

func (r moduleRouter) halt(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id != "" {
		err := r.service.Halt(id)
		if err != nil {
			http.Error(w, "invalid module name", 401)
		}
	}
	w.WriteHeader(200)
}
