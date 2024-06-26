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

type subroutineRouter struct {
	service ports.SubRoutineService
}

func (r *subroutineRouter) RouteInternal(router chi.Router) {
	router.Post("/subroutines/create", r.create)
	router.Post("/subroutines/{id}/run", r.run)
	router.Post("/subroutines/{id}/delete", r.delete)
	router.Post("/subroutines/{id}/update", r.update)
	router.Post("/subroutines/{id}/macros/{macro}/add", r.addMacro)
	router.Post("/subroutines/{id}/macros/{macro}/remove", r.removeMacro)
}

func (r *subroutineRouter) RouteExternal(_ chi.Router) {

}

func NewSubroutineRouter(service ports.SubRoutineService) Routable {
	return &subroutineRouter{
		service: service,
	}
}

func (r *subroutineRouter) removeMacro(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "subroutine key not provided", 401)
		return
	}

	macro := chi.URLParam(req, "macro")
	if macro == "" {
		http.Error(w, "macro key not provided", 401)
		return
	}

	err := r.service.RemoveMacro(id, macro)
	if err != nil {
		http.Error(w, "failed to delete association", 401)
		return
	}
	w.WriteHeader(200)
}

func (r *subroutineRouter) addMacro(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "subroutine key not provided", 401)
		return
	}

	macro := chi.URLParam(req, "macro")
	if macro == "" {
		http.Error(w, "macro key not provided", 401)
		return
	}

	w.WriteHeader(200)
}

func (r *subroutineRouter) update(w http.ResponseWriter, req *http.Request) {
	idNum := chi.URLParam(req, "id")
	if idNum == "" {
		http.Error(w, "subroutine key not provided", 401)
		return
	}
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(req.Body)
	defer req.Body.Close()
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	sr := domain.SubRoutine{}

	err = json.Unmarshal(buf.Bytes(), &sr)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	if idNum != sr.Id {
		http.Error(w, "non-matching delta id", 401)
		return
	}

	id, err := r.service.FindById(idNum)
	if err != nil {
		return
	}

	id.Description = sr.Description
	id.RevertAfter = sr.RevertAfter
	id.Icon = sr.Icon
	id.Group = sr.Group
	id.TriggerId = sr.TriggerId
	id.Macros = sr.Macros

	err = r.service.Update(id)
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

func (r *subroutineRouter) run(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "id")
	if key == "" {
		http.Error(w, "access key not provided", 401)
		return
	}
	err := r.service.Run(key)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
}

func (r *subroutineRouter) delete(w http.ResponseWriter, req *http.Request) {
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

func (r *subroutineRouter) create(w http.ResponseWriter, req *http.Request) {

	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(req.Body)
	defer req.Body.Close()
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	sr := domain.SubRoutine{}

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
