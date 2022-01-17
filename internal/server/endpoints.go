// Copyright (c) 2022 Braden Nicholson

package server

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"udap/internal/auth"
	"udap/internal/bond"
	"udap/internal/cache"
	"udap/internal/controller"
	"udap/internal/log"
	"udap/internal/models"
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
	e.router.Use(auth.VerifyToken())
	// Load JWT Keys
	auth.LoadKeys()

	e.router.Get("/socket/{token}", e.socketAdaptor)
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

func (e *Endpoints) entityBroadcast(ent models.Entity) error {
	response := controller.Response{
		Status:    "success",
		Operation: "entity",
		Body:      ent,
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
		MaxAge:           300, // Maximum value not ignored by any of major browsers
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

	id, err := auth.AuthToken(tokenParam)
	if err != nil {
		log.Err(err)
		return
	}
	ep := e.ctrl.Endpoints.Find(id)
	err = ep.Enroll(c)
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		ep.Connection.Watch()
	}()

	err = e.Metadata()
	if err != nil {
		return
	}

	go func() {
		defer wg.Done()
		for {
			_, out, err := ep.Connection.WS.ReadMessage()
			if err != nil {
				log.Err(err)
				return
			}

			_, err = e.bond.CmdJSON(out)
			if err != nil {
				log.Err(err)
			}
		}
	}()
	wg.Wait()
}

func (e *Endpoints) Broadcast(body any) error {
	for _, s := range e.ctrl.Endpoints.Keys() {
		endpoint := e.ctrl.Endpoints.Find(s)
		if endpoint.Enrolled() {
			endpoint.Connection.Send(body)
		}
	}
	return nil
}

func (e *Endpoints) Run() error {
	log.Log("Endpoints: Listening")

	err := http.ListenAndServe(":3020", e.router)
	if err != nil {
		log.Err(err)
	}
	log.Log("Get: Exiting")
	return nil
}

func (e *Endpoints) registerEndpoint(w http.ResponseWriter, rq *http.Request) {
	key := chi.URLParam(rq, "accessKey")
	endpoint := models.Endpoint{}

	err := store.DB.Model(&models.Endpoint{}).Where("key = ?", key).First(&endpoint).Error
	if err != nil {
	}

	jwt, err := auth.SignUUID(endpoint.Id)
	if err != nil {
		return
	}

	resolve := map[string]interface{}{"token": jwt}

	marshal, err := json.Marshal(resolve)
	if err != nil {
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
	Endpoint  models.Endpoint   `json:"endpoint"`
	Endpoints []models.Endpoint `json:"endpoints"`
	Devices   []models.Device   `json:"devices"`
	Entities  []models.Entity   `json:"entities"`
	Networks  []models.Network  `json:"networks"`
	Logs      []models.Log      `json:"logs"`
	System    System            `json:"system"`
}

func (e *Endpoints) Metadata() error {
	entities, err := e.ctrl.Entities.Compile()
	for _, entity := range entities {
		err = e.entityBroadcast(entity)
		if err != nil {
			return err
		}
	}

	attributes := e.ctrl.Attributes.Compile()
	for _, attribute := range attributes {
		if !e.watching[attribute.Id] {
			cache.WatchFn(attribute.Path(), func(s string) error {
				var attr models.Attribute
				err = json.Unmarshal([]byte(s), &attr)
				if err != nil {
					return err
				}
				err = e.attributeBroadcast(attr)
				if err != nil {
					return err
				}
				return nil
			})
			e.watching[attribute.Id] = true
		}
		err = e.attributeBroadcast(attribute)
		if err != nil {
			return err
		}
	}

	endpoints, err := e.ctrl.Endpoints.Compile()
	if err != nil {
		return err
	}

	devices, err := e.ctrl.Devices.Compile()
	if err != nil {
		return err
	}

	networks, err := e.ctrl.Networks.Compile()
	if err != nil {
		return err
	}

	response := controller.Response{
		Status:    "success",
		Operation: "metadata",
		Body: Metadata{
			Endpoints: endpoints,
			Devices:   devices,
			Networks:  networks,
			System:    SystemInfo,
		},
	}

	err = e.Broadcast(response)
	if err != nil {
		return err
	}

	return nil
}

func (e *Endpoints) Update() error {
	err := e.Metadata()
	if err != nil {
		return err
	}
	return nil
}