// Copyright (c) 2022 Braden Nicholson

package rest

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
	return moduleRouter{
		service: service,
	}
}

func (r moduleRouter) RouteModules(router chi.Router) {
	router.Route("/modules", func(local chi.Router) {
		local.Route("/{name}", func(named chi.Router) {
			named.Post("/build", r.build)
			named.Post("/disable", r.disable)
			named.Post("/enable", r.enable)
			named.Post("/halt", r.halt)
		})
	})
}

func (r moduleRouter) build(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	if name != "" {
		err := r.service.Build(name)
		if err != nil {
			http.Error(w, "invalid module name", 401)
		}
	}
	w.WriteHeader(200)
}

func (r moduleRouter) enable(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	if name != "" {
		err := r.service.Enable(name)
		if err != nil {
			http.Error(w, "invalid module name", 401)
		}
	}
	w.WriteHeader(200)
}

func (r moduleRouter) disable(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	if name != "" {
		err := r.service.Disable(name)
		if err != nil {
			http.Error(w, "invalid module name", 401)
		}
	}
	w.WriteHeader(200)
}

func (r moduleRouter) halt(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	if name != "" {
		err := r.service.Halt(name)
		if err != nil {
			http.Error(w, "invalid module name", 401)
		}
	}
	w.WriteHeader(200)
}
