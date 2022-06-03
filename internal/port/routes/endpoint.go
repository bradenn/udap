// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"net/http"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/platform/jwt"
)

type EndpointRouter interface {
	RouteEndpoints(router chi.Router)
}

type endpointRouter struct {
	service domain.EndpointService
}

func NewEndpointRouter(service domain.EndpointService) EndpointRouter {
	return &endpointRouter{
		service: service,
	}
}

func (r *endpointRouter) RouteEndpoints(router chi.Router) {
	router.Get("/socket/{token}", r.enroll)
	router.Route("/endpoints", func(local chi.Router) {
		local.Get("/register/{key}", r.authenticate)
	})
}

type authenticationResponse struct {
	Token string `json:"token"`
}

func (r *endpointRouter) authenticate(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "key")
	if key == "" {
		http.Error(w, "access key not provided", 401)
		return
	}

	endpoint, err := r.service.FindByKey(key)
	if err != nil {
		http.Error(w, "invalid endpoint name", 401)
		return
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

func (r *endpointRouter) enroll(w http.ResponseWriter, req *http.Request) {
	// Initialize an error to manage returns
	var err error
	// Convert the basic GET request into a WebSocket session
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// Upgrade the https session to a web socket session
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Err(err)
		return
	}
	// find the auth token in the url params
	tokenParam := chi.URLParam(req, "token")
	// Defer the termination of the session to function return
	id, err := jwt.AuthToken(tokenParam)
	if err != nil {
		log.Err(err)
		return
	}

	err = r.service.Enroll(id, conn)
	if err != nil {
		return
	}

	err = r.service.Disconnect(id)
	if err != nil {
		return
	}
}
