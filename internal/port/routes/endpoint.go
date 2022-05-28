// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/core/domain"
	"udap/platform/jwt"
)

type EndpointRouter interface {
	RouteEndpoints(router chi.Router)
}

type endpointRouter struct {
	service domain.EndpointService
}

func NewEndpointRouter(service domain.EndpointService) EndpointRouter {
	return endpointRouter{
		service: service,
	}
}

func (r endpointRouter) RouteEndpoints(router chi.Router) {
	router.Route("/endpoints", func(local chi.Router) {
		local.Post("/authenticate/{key}", r.authenticate)
	})
}

type authenticationResponse struct {
	Token string `json:"token"`
}

func (r endpointRouter) authenticate(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "key")
	if key == "" {
		http.Error(w, "access key not provided", 401)
	}

	endpoint, err := r.service.FindByKey(key)
	if err != nil {
		http.Error(w, "invalid endpoint name", 401)
	}

	token, err := jwt.SignUUID(endpoint.Id)
	if err != nil {
		http.Error(w, "Failed to generate JWT.", 500)
		return
	}

	resolve := authenticationResponse{}
	resolve.Token = token

	marshal, err := json.Marshal(resolve)
	if err != nil {
		http.Error(w, "Failed to generate json...", 500)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		return
	}
	w.WriteHeader(200)
}
