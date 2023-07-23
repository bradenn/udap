// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"net/http"
	"udap/internal/core/domain"
	"udap/internal/core/ports"
	"udap/internal/log"
	"udap/platform/jwt"
)

type endpointRouter struct {
	service ports.EndpointService
}

func NewEndpointRouter(service ports.EndpointService) Routable {
	return &endpointRouter{
		service: service,
	}
}

func (r *endpointRouter) RouteExternal(router chi.Router) {
	router.Get("/endpoints/register/{key}", r.authenticate)
	router.Post("/endpoints/create", r.create)
	router.Post("/endpoints/{id}/push", r.registerPush)
	router.Get("/socket/{token}", r.enroll)
}

func (r *endpointRouter) RouteInternal(router chi.Router) {

}

func (r *endpointRouter) registerPush(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	key := chi.URLParam(req, "id")
	if key == "" {
		http.Error(w, "access key not provided", 401)
		return
	}

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(w, "could not parse endpoint", 400)
		return
	}

	err = r.service.RegisterPush(key, buf.String())
	if err != nil {
		http.Error(w, "endpoint creation failed", 400)
	}
}
func (r *endpointRouter) create(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(w, "could not parse endpoint", 400)
		return
	}
	ref := domain.Endpoint{}
	err = json.Unmarshal(buf.Bytes(), &ref)
	if err != nil {
		http.Error(w, "could not parse endpoint", 400)
		return
	}

	product := domain.NewEndpoint(ref.Name, ref.Type)

	err = r.service.Create(product)
	if err != nil {
		http.Error(w, "endpoint creation failed", 400)
	}
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
		log.Err(err)
	}

	done := make(chan bool)
	go func() {
		_, _, err = conn.ReadMessage()
		if err != nil {
			done <- true
		}
	}()

	defer func() {
		_ = conn.Close()
		err = r.service.Unenroll(id)
		if err != nil {
			log.Err(err)
		}
	}()

	<-done

}
