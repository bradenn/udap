// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/core/domain"
	"udap/internal/log"
)

type moduleRouter struct {
	service domain.ModuleService
}

func NewModuleRouter(service domain.ModuleService) Routable {
	return &moduleRouter{
		service: service,
	}
}

func (r *moduleRouter) RouteInternal(router chi.Router) {
	router.Route("/modules/{id}", func(local chi.Router) {
		local.Post("/reload", r.reload)
		local.Post("/build", r.build)
		local.Post("/disable", r.disable)
		local.Post("/enable", r.enable)
		local.Post("/halt", r.halt)
	})
}

func (r *moduleRouter) RouteExternal(_ chi.Router) {

}

func (r *moduleRouter) reload(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id != "" {
		err := r.service.Reload(id)
		if err != nil {
			log.Err(err)
			http.Error(w, "an error occured: ", 401)
			return
		}
	}
	w.WriteHeader(200)
}

func (r *moduleRouter) build(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id != "" {
		mod, err := r.service.FindByName(id)
		if err != nil {
			return
		}
		err = r.service.Build(mod.Id)
		if err != nil {
			http.Error(w, "invalid module name", 401)
		}
	}
	w.WriteHeader(200)
}

func (r *moduleRouter) enable(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id != "" {
		err := r.service.Enable(id)
		if err != nil {
			http.Error(w, "invalid module name", 401)
		}
	}
	w.WriteHeader(200)
}

func (r *moduleRouter) disable(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id != "" {
		err := r.service.Disable(id)
		if err != nil {
			http.Error(w, "invalid module name", 401)
		}
	}
	w.WriteHeader(200)
}

func (r *moduleRouter) halt(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id != "" {
		err := r.service.Halt(id)
		if err != nil {
			http.Error(w, "invalid module name", 401)
		}
	}
	w.WriteHeader(200)
}
