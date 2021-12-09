// Copyright (c) 2021 Braden Nicholson

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"udap/internal/cache"
	"udap/internal/log"
	"udap/internal/models"
	"udap/pkg/plugin"
)

// Runtime manages the event-loop for all instances, as well as the websocket connections between UDAP and endpoints.
type Runtime struct {
	modules   Modules   // Key: "name"
	entities  Entities  // Key: "module.name"
	endpoints Endpoints // Key: "uuid"

	eventHandler chan plugin.Event
}

// Dependency is the level at which this service needs to run
func (r *Runtime) Dependency() (level int) {
	return 1
}

func (r *Runtime) Name() (name string) {
	return "runtime"
}

func (r *Runtime) Load() (err error) {
	// Initialize the module and entity poly buffers
	r.modules = Modules{}
	r.entities = Entities{}
	r.endpoints = Endpoints{}

	// Initialize a bounded buffer for accepting module events
	r.eventHandler = make(chan plugin.Event, 1)
	// Launch the event resolution thread
	go func() {
		for event := range r.eventHandler {
			switch event.Type {
			case "entity":
				err = r.entityEvent(event)
				if err != nil {
					return
				}
			}
		}
		close(r.eventHandler)
	}()

	err = r.buildModules()
	if err != nil {
		return err
	}
	return nil
}

// Update is called when the runtime is to begin accepting traffic
func (r *Runtime) Update() {
	// The state of each entity will be polled every 2 seconds
	for _, s := range r.entities.Keys() {
		// Get the entity from the mutex map
		e := r.entities.Find(s)
		// Attempt to poll it
		err := e.Poll()
		if err != nil {
			return
		}
		// Save the state to the cache
		err = cache.PutLn(e.State, e.Path(), "state")
		if err != nil {
			return
		}
	}
}

// Run is called when the runtime is to begin accepting traffic
func (r *Runtime) Run(ctx context.Context) error {
	r.endpoints.FetchAll()
	// Attempt to get the server reference from context
	srv := ctx.Value("server").(*Server)
	// Use the server from context to handle endpoint socket connections
	srv.Get("/socket/{token}", r.HandleSockets)
	// Attempt to load the modules in the directory 'modules'
	err := r.loadModulesDir("modules")
	if err != nil {
		return err
	}
	// Create a wait group so all plugins can init at the same time
	wg := sync.WaitGroup{}
	wg.Add(len(r.modules.Keys()))
	// Run the full lifecycle of all plugins
	for _, s := range r.modules.Keys() {
		// Find the UdapPlugin module
		pl := r.modules.Find(s)
		// Run a go function to create a new thread
		go func(p plugin.UdapPlugin) {
			// Defer the wait group to complete at the end
			defer wg.Done()
			// Run module setup
			_, err = p.Setup()
			if err != nil {
				log.Err(err)
				return
			}
			// Attempt to connect to the module
			_, err = p.Connect(&r.eventHandler)
			if err != nil {
				log.Err(err)
				return
			}
			// Attempt to run the module
			err = p.Run()
			if err != nil {
				log.Err(err)
				return
			}
		}(pl)
	}
	// Wait for all modules to exit
	wg.Wait()
	return nil
}

func (r *Runtime) entityEvent(event plugin.Event) (err error) {
	entity := event.Body.(*models.Entity)
	err = entity.Emplace()
	if err != nil {
		return err
	}
	r.entities.register(entity.Path(), entity)
	log.Log("Entity '%s' loaded", entity.Name)
	return nil
}

// buildModules locates and initializes modules in the plugin directory
func (r *Runtime) buildModules() error {
	// Try to load modules from the plugin folders
	wg := sync.WaitGroup{}
	err := r.buildModuleDir("modules", &wg)
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

// buildModuleDir builds all potential modules in a directory
func (r *Runtime) buildModuleDir(dir string, wg *sync.WaitGroup) error {
	// Format the pattern for glob search
	pattern := fmt.Sprintf("./plugins/%s/*/*.go", dir)
	// Run the search for go files
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	// Add all the potential files from the search
	wg.Add(len(files))
	// Launch a go func to build each one
	for _, p := range files {
		// Run the function for this file
		go func(path string) {
			if err := r.buildFromSource(path, wg); err != nil {
				// If an error occurs, print it to console
				log.ErrF(err, "failed to build module candidate '%s'", path)
			}
		}(p)
	}
	return nil
}

// buildFromSource will build an eligible plugin from sources if applicable
func (r *Runtime) buildFromSource(path string, wg *sync.WaitGroup) error {
	// Create output file by modifying input file extension
	out := strings.Replace(path, ".go", ".so", 1)
	// Create a timeout to prevent modules from taking too long to build
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*4)
	// Cancel the timeout of it exits before the timeout is up
	defer cancelFunc()
	// Get the go executable from the environment
	goExec := os.Getenv("goExec")
	// Prepare the command arguments
	args := []string{"build", "-v", "-buildmode=plugin", "-o", out, path}
	// Initialize the command structure
	cmd := exec.CommandContext(timeout, goExec, args...)
	// Run and get the stdout and stderr from the output
	err := cmd.Run()
	if err != nil {
		return err
	}
	// Mark this plugin as complete
	wg.Done()
	return nil
}

// loadModulesDir attempts to load each module
func (r *Runtime) loadModulesDir(dir string) error {
	path := fmt.Sprintf("./plugins/%s/*/*.so", dir)
	files, err := filepath.Glob(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		p, err := plugin.Load(file)
		if err != nil {
			log.Err(err)
			continue
		}
		name := strings.Replace(filepath.Base(file), ".so", "", 1)

		r.modules.Set(name, p)

	}
	return nil
}

// loadMutable attempts to load each module
func (r *Runtime) loadMutable() error {
	return nil
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

// HandleSockets upgrades a http request and initiates a WebSocket session
func (r *Runtime) HandleSockets(w http.ResponseWriter, req *http.Request) {
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
	sender := r.endpoints.Find(id)
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

func (r *Runtime) routeRequests(request RequestS) (err error) {
	switch t := request.Target; t {
	case "endpoint":
		err = r.EndpointRequest(request)
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
	Name   string          `json:"name"`   // EntityId
	Module string          `json:"module"` // EntityId
	State  json.RawMessage `json:"state"`  // EntityId

}

func (r *Runtime) EntityRequest(request RequestS) (err error) {
	var payload EntityState
	err = json.Unmarshal(request.Body, &payload)
	if err != nil {
		return err
	}
	entity := models.Entity{}
	entity.Name = payload.Name
	entity.Module = payload.Module
	switch t := request.Operation; t {
	case "state":
		log.Log("%s", payload.State)
		e := r.entities.Find(entity.Path())
		err = e.Push(models.State(payload.State))
		if err != nil {
			return err
		}
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
	// Find all active entities
	var entities []models.Entity
	for _, entity := range r.entities.Keys() {
		m := r.entities.Find(entity)
		if m.Live() {
			err := m.Poll()
			if err != nil {
				return err
			}
		}
		entities = append(entities, *m)
	}
	// Create and populate the response body
	response.Body = Object{
		"endpoint": e,
		"entities": entities,
	}
	// Convert the struct into json
	marshal, err := json.Marshal(response)
	if err != nil {
		return err
	}
	// Send the JSON via the endpoints Conn pointer
	err = e.Conn.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		return err
	}
	// Return no errors
	return nil
}
