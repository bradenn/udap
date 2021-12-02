// Copyright (c) 2021 Braden Nicholson

package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"net/http"
	"udap/internal/log"
	"udap/internal/models"
	"udap/pkg/plugin"
)

// Target: Entity
// Id: EntityId
// Payload to send: {StateChangePayload}

// Runtime manages the event-loop for all instances, as well as the websocket connections between UDAP and endpoints.
type Runtime struct {
	resolver chan plugin.Event
	plugins  map[string]*plugin.Plugin
	modules  map[string]*models.Module
}

// Dependency is the level at which this service needs to run
func (r *Runtime) Dependency() (level int) {
	return 3
}

func (r *Runtime) Name() (name string) {
	return "runtime"
}

func (r *Runtime) Load() (err error) {
	// The updater channel is called by modules, we initialize it to queue 16 at a time
	resolver := make(chan plugin.Event, 16)

	r.resolver = resolver
	r.plugins = map[string]*plugin.Plugin{}
	r.modules = map[string]*models.Module{}

	// Discover any unloaded modules
	r.discoverModules()
	// return the runtime
	return nil
}

// Run is called when the runtime is to begin accepting traffic
func (r *Runtime) Run(interface{}) error {
	// Set async listener for the updating from modules and daemons
	server.router.Get("/socket/{token}", r.HandleSockets)
	return nil
}

func (r *Runtime) Cleanup() (err error) {
	return nil
}

// discoverModules locates and initializes modules in the plugin directory
func (r *Runtime) discoverModules() {
	log.Sherlock("Discovering modules...")
	// Try to load modules from the plugin folders
	plugins := plugin.LoadAll("./plugins/modules")
	// if err != nil {
	// 	log.Err(err)
	// }
	for path, p := range plugins {
		err := r.initializeModule(path, p)
		if err != nil {
			log.Err(err)
			return
		}
	}
}

// initializeModule loads a plugin in as a module
func (r *Runtime) initializeModule(path string, p plugin.Plugin) error {
	_, err := p.Startup()
	if err != nil {
		return err
	}
	m, err := models.FromPath(path, p)
	if err != nil {
		return err
	}
	// Provide the resolver channel for the plugin to talk with this module
	p.Connect(r.resolver)
	if r.modules[m.Id] != nil {
		return nil
	}
	r.modules[m.Id] = &m
	log.Sherlock("Module '%s' v%s initialized.", m.Name, m.Version)
	// Return no errors
	return nil
}

func (r *Runtime) Resolve() {

	// close(r.resolver)
}

type Update struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Data string `json:"data"`
	Id   string `json:"id"`
}

type RequestS struct {
	Target    string          `json:"target"`
	Operation string          `json:"operation"`
	Body      json.RawMessage `json:"body"`
	Sender    *models.Endpoint
}

type Response struct {
	Status    string                 `json:"status"`
	Operation string                 `json:"operation"`
	Body      map[string]interface{} `json:"body"`
}

// Load modules and instances

// HandleSockets upgrades a http request and initiates a WebSocket session
func (r *Runtime) HandleSockets(w http.ResponseWriter, req *http.Request) {
	// Initialize an error to manage returns
	var err error
	// Convert the basic GET request into a WebSocket session
	c, err := Upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Err(err)
		return
	}
	// Find the auth token in the url params
	tokenParam := chi.URLParam(req, "token")
	// Defer the termination of the session to function return
	defer func(c *websocket.Conn) {
		_ = c.Close()
	}(c)

	// Authenticate the endpoint by its JWT token
	id, err := AuthToken(tokenParam)
	if err != nil {
		log.Err(err)
		return
	}
	// Attempt to find the endpoint based on its UUID
	sender := &models.Endpoint{}
	sender.Id = id
	err = sender.Fetch()
	if err != nil {
		return
	}
	sender.Conn = c
	// RequestS event loop, runs every time a message is received
	for {
		// Initialize the request structure
		request := RequestS{
			Sender: sender,
		}
		// Wait for a request from the client, ReadJSON will halt this loop indefinitely if needed
		err = c.ReadJSON(&request)
		if err != nil {
			r.sendError(request, err)
			return
		}
		// Determine the request's intention before assigning its proper function
		err = r.routeRequests(request)
		if err != nil {
			r.sendError(request, err)
			return
		}

		err = r.Metadata(sender)
		if err != nil {
			r.sendError(request, err)
			return
		}
	}
}

func (r *Runtime) sendError(request RequestS, err error) {
	c := request.Sender.Conn
	log.ErrF(err, "Endpoint '%s' disconnected due to an error", request.Sender.Name)
	// Initialize error message response
	response := Response{
		Status:    "error",
		Operation: request.Operation,
		Body:      map[string]interface{}{"message": err.Error()},
	}
	// Write the struct to the connection
	err = c.WriteJSON(response)
	if err != nil {
		log.Err(err)
		return
	}
	// Attempt to close the connection
	err = c.Close()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// path : [namespace].[identifier].[inquiry]
// { Action: "exit", Args: []string{"anal"} } => instance.appls23ddjskj
// P("module.%s", moduleId)
// { Command:}

func (r *Runtime) routeRequests(request RequestS) (err error) {
	switch t := request.Target; t {
	case "endpoint":
		err = r.EndpointRequest(request)
		if err != nil {
			return err
		}
	case "instance":
		err = r.InstanceRequest(request)
		if err != nil {
			return err
		}
	case "entity":
		err = r.EntityRequest(request)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("the target '%s' is unknown", t)
	}
	return nil
}

type EntityState struct {
	Id    string `json:"id"`    // EntityId
	State string `json:"state"` // EntityId

}

func (r *Runtime) EntityRequest(request RequestS) (err error) {
	var payload EntityState
	err = json.Unmarshal(request.Body, &payload)
	if err != nil {
		return err
	}
	entity := models.Entity{}
	entity.Id = payload.Id
	err = entity.Fetch()
	if err != nil {
		return err
	}
	switch t := request.Operation; t {
	case "state":
		entity.SetState(payload.State)
	}
	err = r.Metadata(request.Sender)
	if err != nil {
		return err
	}
	return nil
}

func (r *Runtime) EndpointRequest(request RequestS) (err error) {
	var reqBody models.RequestBody
	err = json.Unmarshal(request.Body, &reqBody)
	if err != nil {
		return err
	}
	switch t := request.Operation; t {
	case "metadata":
		err = r.Metadata(request.Sender)
		if err != nil {
			return err
		}
	case "rename":
		// err = request.Sender.Rename(reqBody)
		if err != nil {
			return err
		}
	case "enroll":
		err = request.Sender.Enroll()
		if err != nil {
			return err
		}
	case "unenroll":
		err = request.Sender.Unenroll()
		if err != nil {
			return err
		}
	case "grant":
		err = request.Sender.Grant(reqBody)
		if err != nil {
			return err
		}
	case "revoke":
		err = request.Sender.Revoke(reqBody)
		if err != nil {
			return err
		}
	case "subscribe":
		err = request.Sender.Subscribe(reqBody)
		if err != nil {
			return err
		}
	case "unsubscribe":
		err = request.Sender.Unsubscribe(reqBody)
		if err != nil {
			return err
		}
	}
	err = r.Metadata(request.Sender)
	if err != nil {
		return err
	}
	return nil
}

type IdentifierBody struct {
	Id string `json:"id"`
}

func (r *Runtime) InstanceRequest(request RequestS) (err error) {
	body := IdentifierBody{}
	i := &models.Instance{}

	err = json.Unmarshal(request.Body, &body)
	if err != nil {
		return err
	}

	if body.Id == "" {
		switch t := request.Operation; t {
		case "create":
			err = i.Create(request.Body)
			if err != nil {
				return err
			}
			// TODO Make it grant new ones to the endpoint it was made on.
		default:
			return fmt.Errorf("invalid operation '%s'; you may have forgotten to provice an Id", t)
		}
		return nil
	}
	i, err = models.GetInstance(body.Id)
	if err != nil {
		return err
	}
	switch t := request.Operation; t {
	case "modify":
		err = i.Modify(request.Body)
		if err != nil {
			return err
		}
	case "run":
		err = i.Run(string(request.Body))
		if err != nil {
			return err
		}
	case "reset":
		err = i.Reset()
		if err != nil {
			return err
		}
	case "delete":
		err = i.Reset()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid operation '%s'", t)
	}

	return nil
}

type Object map[string]interface{}

// Metadata sends metadata back to an enrolled endpoint
func (r *Runtime) Metadata(e *models.Endpoint) error {
	// Initialize a response struct
	response := Response{
		Status:    "success",
		Operation: "metadata",
		Body:      Object{},
	}
	// Collect loaded modules
	var modules []models.Module
	for _, module := range r.modules {
		modules = append(modules, *module)
	}
	// Find all active entities
	entities, err := models.GetEntities()
	if err != nil {
		log.Err(err)
	}
	// Get all active instances
	instances, err := models.GetInstances()
	if err != nil {
		log.Err(err)
	}
	// Create and populate the response body
	response.Body = Object{
		"endpoint":  e,
		"instances": instances,
		"entities":  entities,
		"modules":   modules,
	}
	// Convert the struct into json
	marshal, err := json.Marshal(response)
	if err != nil {
		log.Err(err)
	}
	// Send the JSON via the endpoints Conn pointer
	err = e.Conn.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		log.Error("write:", err)
		log.Err(err)
	}
	// Return no errors
	return nil
}

// Push sends data to all enrolled endpoints
func Push(e models.Endpoint, instances map[string]interface{}) {

	response := Response{
		Status:    "success",
		Operation: "update",
		Body:      map[string]interface{}{},
	}

	// for _, instance := range e.Subscriptions {
	// 	response.Body[instance.UUID()] = instances[instance.UUID()]
	// }

	marshal, err := json.Marshal(response)
	if err != nil {
		return
	}

	err = e.Conn.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		log.Error("write:", err)
		return
	}

}
