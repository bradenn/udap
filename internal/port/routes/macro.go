// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/core/ports"
)

type macroRouter struct {
	service ports.MacroService
}

func NewMacroRouter(service ports.MacroService) Routable {
	return macroRouter{
		service: service,
	}
}

func (r macroRouter) RouteInternal(router chi.Router) {
	router.Route("/macros", func(local chi.Router) {
		local.Post("/{id}/run", r.run)
	})
}

func (r macroRouter) RouteExternal(_ chi.Router) {

}

func (r macroRouter) run(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "id")
	if key == "" {
		http.Error(w, "access key not provided", 401)
		return
	}

	err := r.service.Run(key)
	if err != nil {
		http.Error(w, "could not run macro", 500)
		return
	}

	w.WriteHeader(200)
}
