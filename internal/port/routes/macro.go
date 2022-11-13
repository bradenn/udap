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
	router.Post("/macros/create", r.create)
	router.Post("/macros/{id}/delete", r.delete)
	router.Post("/macros/{id}/run", r.run)
	router.Post("/macros/{id}/update", r.update)
}

func (r macroRouter) RouteExternal(_ chi.Router) {

}

func (r macroRouter) delete(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "id")
	if key == "" {
		http.Error(w, "access key not provided", 401)
		return
	}

	err := r.service.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
}

func (r macroRouter) update(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "id")
	if key == "" {
		http.Error(w, "access key not provided", 401)
		return
	}
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

	byId, err := r.service.FindById(key)
	if err != nil {
		return
	}

	byId.Name = sr.Name
	byId.Description = sr.Description
	byId.ZoneId = sr.ZoneId
	byId.Type = sr.Type
	byId.Value = sr.Value

	err = r.service.Update(&sr)
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
