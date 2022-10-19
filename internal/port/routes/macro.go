// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"udap/internal/core/domain"
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
		local.Post("/create", r.create)
		local.Post("/{id}/run", r.run)
	})
}

func (r macroRouter) RouteExternal(_ chi.Router) {

}

func (r macroRouter) create(w http.ResponseWriter, req *http.Request) {

	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(req.Body)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(req.Body)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	sr := domain.Macro{}

	err = json.Unmarshal(buf.Bytes(), &sr)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	err = r.service.Create(&sr)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	_, err = w.Write([]byte("OK"))
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(200)
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
