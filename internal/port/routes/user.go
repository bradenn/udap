// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/core/domain"
)

type userRouter struct {
	service domain.UserService
}

func NewUserRouter(service domain.UserService) Routable {
	return userRouter{
		service: service,
	}
}

func (r userRouter) RouteInternal(router chi.Router) {
	router.Route("/users", func(local chi.Router) {
		local.Post("/register", r.register)
		local.Post("/authenticate", r.authenticate)
	})
}

func (r userRouter) RouteExternal(_ chi.Router) {

}

func (r userRouter) authenticate(w http.ResponseWriter, req *http.Request) {

	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return
	}
	ref := domain.User{}
	err = json.Unmarshal(buf.Bytes(), &ref)
	if err != nil {
		http.Error(w, "could not parse user", 400)
		return
	}

	err = r.service.Authenticate(&ref)
	if err != nil {
		http.Error(w, "failed to create user", 400)
		return
	}
	w.WriteHeader(200)
}

func (r userRouter) register(w http.ResponseWriter, req *http.Request) {

	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return
	}
	ref := domain.User{}
	err = json.Unmarshal(buf.Bytes(), &ref)
	if err != nil {
		http.Error(w, "could not parse user", 400)
		return
	}

	err = r.service.Register(&ref)
	if err != nil {
		http.Error(w, "failed to create user", 400)
		return
	}
	w.WriteHeader(200)
}
