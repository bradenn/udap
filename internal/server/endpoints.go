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
	"sync"
	"udap/internal/auth"
	"udap/internal/bond"
	"udap/internal/controller"
	"udap/internal/log"
	"udap/internal/models"
	"udap/internal/store"
)

type Endpoints struct {
	bond        *bond.Bond
	router      chi.Router
	connections sync.Map
}

func (e *Endpoints) Setup(bond *bond.Bond) error {
	e.bond = bond
	e.router = chi.NewRouter()

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
	// Find the auth token in the url params
	tokenParam := chi.URLParam(req, "token")
	// Defer the termination of the session to function return

	id, err := auth.AuthToken(tokenParam)
	if err != nil {
		log.Err(err)
		return
	}
	cn := NewConnection(c)
	e.connections.Store(id, cn)

	c.SetCloseHandler(func(code int, text string) error {
		andDelete, loaded := e.connections.LoadAndDelete(id)
		if loaded {
			close(andDelete.(*Connection).edit)
		}
		return nil
	})

	err = e.HandleRequest(cn)
	if err != nil {
		log.Err(err)
	}

}

func (e *Endpoints) HandleRequest(connection *Connection) error {
	for {
		_, out, err := connection.WS.ReadMessage()
		if err != nil {
			return err
		}

		t, err := e.bond.CmdJSON(out)
		if err != nil {
			return err
		}
		fmt.Println(t)

	}
}

func (e *Endpoints) Run() error {
	log.Log("Endpoints: Listening")
	err := http.ListenAndServe(":3020", e.router)
	if err != nil {
		log.Err(err)
	}
	log.Log("Remote: Exiting")
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

// Metadata sends metadata back to an enrolled endpoint
func (e *Endpoints) metadata(msg bond.Msg) error {
	// err := e.Metadata(msg.Id)
	// if err != nil {
	// 	return err
	// }
	return nil
}

type Connection struct {
	WS   *websocket.Conn
	edit chan any
}

func (c *Connection) Send(body any) {
	c.edit <- body
}

func NewConnection(ws *websocket.Conn) *Connection {
	ch := make(chan any, 16)
	c := &Connection{
		WS:   ws,
		edit: ch,
	}
	go c.Watch()

	return c
}

func (c *Connection) Watch() {
	for msg := range c.edit {
		if c.WS == nil {
			continue
		}
		err := c.WS.WriteJSON(msg)
		if err != nil {
			log.Err(err)
		}
	}
	c.WS = nil
}

func (e *Endpoints) Metadata(id string, connection *Connection) error {

	entities, err := e.bond.Send("entity", "compile", nil)
	if err != nil {
		return err
	}

	endpoints, err := e.bond.Send("endpoint", "compile", nil)
	if err != nil {
		return err
	}

	devices, err := e.bond.Send("device", "compile", nil)
	if err != nil {
		return err
	}

	networks, err := e.bond.Send("network", "compile", nil)
	if err != nil {
		return err
	}

	var logs []models.Log
	err = store.DB.Model(&models.Log{}).Order("created_at desc").Limit(40).Find(&logs).Error
	if err != nil {
		return err
	}

	mid := models.Endpoint{}
	mid.Id = id

	response := controller.Response{
		Status:    "success",
		Operation: "metadata",
		Body: controller.Metadata{
			Logs:      logs,
			Entities:  entities.([]models.Entity),
			Devices:   devices.([]models.Device),
			Networks:  networks.([]models.Network),
			Endpoints: endpoints.([]models.Endpoint),
			Endpoint:  mid,
		},
	}

	connection.Send(response)

	return nil
}

func (e *Endpoints) Update() error {
	e.connections.Range(func(key, value any) bool {
		err := e.Metadata(key.(string), value.(*Connection))
		if err != nil {
			log.Err(err)
		}
		return true
	})
	return nil
}
