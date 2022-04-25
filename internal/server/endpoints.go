// Copyright (c) 2022 Braden Nicholson

package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"sync"
	"udap/internal/bond"
	"udap/internal/controller"
	"udap/internal/log"
	"udap/internal/models"
	"udap/internal/pulse"
	"udap/internal/store"
)

type Endpoints struct {
	bond        *bond.Bond
	router      chi.Router
	connections sync.Map
	watching    map[string]bool
	ctrl        *controller.Controller
}

func (e *Endpoints) Name() string {
	return "endpoints"
}

func (e *Endpoints) Setup(ctrl *controller.Controller, bond *bond.Bond) error {
	e.ctrl = ctrl
	e.bond = bond
	e.router = chi.NewRouter()
	e.watching = map[string]bool{}
	e.router.Use(middleware.Recoverer)
	// Custom Middleware
	e.router.Use(corsHeaders())
	// Status Middleware
	e.router.Use(middleware.Heartbeat("/status"))
	// Seek, verify and validate JWT tokens
	e.router.Use(verifyToken())
	// Load JWT Keys
	loadKeys()
	// Route the websocket listening endpoint
	e.router.Get("/socket/{token}", e.socketAdaptor)
	// Route the endpoint registration
	e.router.Get("/endpoints/register/{accessKey}", e.registerEndpoint)
	return nil
}

func (e *Endpoints) attributeBroadcast(ent models.Attribute) error {
	response := controller.Response{
		Status:    "success",
		Operation: "attribute",
		Body:      ent,
	}

	err := e.Broadcast(response)
	if err != nil {
		return err
	}
	return nil
}

func (e *Endpoints) reactive(operation string) func(interface{}) error {
	return func(a interface{}) error {
		response := controller.Response{
			Status:    "success",
			Operation: operation,
			Body:      a,
		}

		err := e.Broadcast(response)
		if err != nil {
			return err
		}

		return nil
	}

}

func (e *Endpoints) itemBroadcast(operation string, body interface{}) error {
	response := controller.Response{
		Status:    "success",
		Operation: operation,
		Body:      body,
	}

	err := e.Broadcast(response)
	if err != nil {
		return err
	}
	return nil
}

func corsHeaders() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Bond"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by interface{} of major browsers
	})
}

func (e *Endpoints) socketAdaptor(w http.ResponseWriter, req *http.Request) {
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
	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Err(err)
		return
	}
	// find the auth token in the url params
	tokenParam := chi.URLParam(req, "token")
	// Defer the termination of the session to function return

	id, err := authToken(tokenParam)
	if err != nil {
		log.Err(err)
		return
	}
	ep := e.ctrl.Endpoints.Find(id)
	err = ep.Enroll(c)
	if err != nil {
		log.Err(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		ep.Connection.Watch()
	}()

	err = e.ctrl.Entities.EmitAll()
	if err != nil {
		return
	}

	err = e.ctrl.Attributes.EmitAll()
	if err != nil {
		return
	}

	err = e.ctrl.Networks.EmitAll()
	if err != nil {
		return
	}

	err = e.ctrl.Devices.EmitAll()
	if err != nil {
		return
	}

	err = e.ctrl.Endpoints.EmitAll()
	if err != nil {
		return
	}

	go func() {
		defer wg.Done()
		for {
			_, out, err := ep.Connection.WS.ReadMessage()
			if err != nil {
				return
			}

			_, err = e.bond.CmdJSON(out)
			if err != nil {
				err = e.sendError(id, err)
				if err != nil {
					return
				}
			}
		}
	}()
	wg.Wait()
}

func (e *Endpoints) sendError(id string, body interface{}) error {
	response := controller.Response{
		Status:    "error",
		Operation: "error",
		Body:      body,
	}
	endpoint := e.ctrl.Endpoints.Find(id)
	if endpoint.Enrolled() {
		endpoint.Connection.Send(response)
	}
	return nil
}

func (e *Endpoints) Broadcast(body interface{}) error {
	for _, s := range e.ctrl.Endpoints.Keys() {
		endpoint := e.ctrl.Endpoints.Find(s)
		if endpoint.Enrolled() {
			endpoint.Connection.Send(body)
		}
	}
	return nil
}

func (e *Endpoints) Run() error {
	port := os.Getenv("hostPort")

	e.ctrl.Devices.Watch(e.reactive("device"))
	e.ctrl.Entities.Watch(e.reactive("entity"))
	e.ctrl.Attributes.Watch(e.reactive("attribute"))
	e.ctrl.Endpoints.Watch(e.reactive("endpoint"))
	e.ctrl.Networks.Watch(e.reactive("network"))

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), e.router)
	if err != nil {
		log.Err(err)
	}
	return nil
}

func (e *Endpoints) registerEndpoint(w http.ResponseWriter, rq *http.Request) {
	key := chi.URLParam(rq, "accessKey")
	ep := models.Endpoint{}

	err := store.DB.Model(&models.Endpoint{}).Where("key = ?", key).First(&ep).Error
	if err != nil {
		http.Error(w, "Invalid security code.", 401)
		return
	}

	jwt, err := signUUID(ep.Id)
	if err != nil {
		http.Error(w, "Failed to generate JWT.", 500)
		return
	}

	resolve := map[string]interface{}{"token": jwt}

	marshal, err := json.Marshal(resolve)
	if err != nil {
		http.Error(w, "Failed to generate json...", 500)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		return
	}
}

type Identifier struct {
	Id string `json:"id"`
}

type Metadata struct {
	System System `json:"system"`
}

func (e *Endpoints) Timings() error {
	timings := pulse.Timings.Timings()
	for _, timing := range timings {
		err := e.itemBroadcast("timing", timing)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Endpoints) Metadata() error {

	response := controller.Response{
		Status:    "success",
		Operation: "metadata",
		Body: Metadata{
			System: SystemInfo,
		},
	}

	err := e.Broadcast(response)
	if err != nil {
		return err
	}

	return nil
}

func (e *Endpoints) Update() error {
	pulse.Fixed(500)
	defer pulse.End()
	err := e.Metadata()
	if err != nil {
		return err
	}
	return nil

}
