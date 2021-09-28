package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"udap/logger"
	"udap/server"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Runtime struct {
	endpoints map[string]*Session
	server    server.Server
	instances map[string]*Instance
}

type Session struct {
	Connection *websocket.Conn
	Instances  []string
}

type Error struct {
	Error interface{} `json:"error"`
}

func NewError(err string) Error {
	return Error{Error: err}
}

type Request struct {
	Token   string          `json:"token"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EnrollmentPayload struct {
	Instances []string `json:"instances"`
}

type ActionPayload struct {
	Instance string `json:"instances"`
	Action   string `json:"action"`
}

type ConfigurationPayload struct {
	Instance string `json:"instances"`
	Action   string `json:"action"`
}

// Enroll adds an endpoint's websocket connection to the active update queue
func (r *Runtime) Enroll(endpoint Endpoint, conn *websocket.Conn, payload json.RawMessage) error {
	// Initialize a structure to contain the Enrollment instances
	data := EnrollmentPayload{}
	// Attempt to parse the payload into the structure
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return err
	}
	// Generate a session struct, containing the websocket connection, and the requested watchlist
	session := Session{
		Connection: conn,
		Instances:  data.Instances,
	}
	// Log the new state
	logger.Info("Activated endpoint '%s'", endpoint.Name)
	// Set a close handler to deactivate the endpoint when the connection is severed
	session.Connection.SetCloseHandler(r.socketCloseHandler(endpoint))
	// Enroll the endpoint into the runtime
	r.endpoints[endpoint.Id.String()] = &session
	// Push a welcome update
	r.UpdateOne(endpoint.Id.String(), &session)
	return nil
}

// Enrolled returns true if the provided endpoint's UUID is in enrolled
func (r *Runtime) Enrolled(endpoint Endpoint) bool {
	return r.endpoints[endpoint.Id.String()] != nil
}

// Unenroll removes an endpoint from the active update queue
func (r *Runtime) Unenroll(endpoint Endpoint) {
	// Handle the base case
	if r.endpoints[endpoint.Id.String()] == nil {
		return
	}
	// Remove the connection reference from the endpoint map
	delete(r.endpoints, endpoint.Id.String())
	// Log the deactivation
	logger.Info("Deactivated endpoint '%s'", endpoint.Name)
}

// EndpointHandler upgrades a http request and initiates a websocket session
func (r *Runtime) EndpointHandler(w http.ResponseWriter, req *http.Request) {
	// Convert the basic GET request into a Websocket session
	c, err := upgrader.Upgrade(w, req, nil)
	// Defer the termination of the session to function return
	defer func(c *websocket.Conn) {
		_ = c.Close()
	}(c)
	// Ignore malformed requests
	if err != nil {
		return
	}
	// Request event loop, runs every time a message is received
	for {
		// Initialize the request structure
		request := Request{}
		// Wait for a request from the client, ReadJSON will halt this loop indefinitely if needed
		err = c.ReadJSON(&request)
		// Reject malformed requests silently
		if err != nil {
			break
		}
		// Authenticate the endpoint by its JWT token
		id, err := server.AuthToken(request.Token)
		// Aggressively notify the client of their critical mistake
		if err != nil {
			_ = c.WriteJSON(NewError(err.Error()))
			_ = c.Close()
			break
		}
		// Attempt to find the endpoint based on its UUID
		endpoint, err := r.findEndpoint(id)
		if err != nil {
			_ = c.WriteJSON(NewError(err.Error()))
			_ = c.Close()
			break
		}
		// Determine the request's intention before assigning its proper function
		switch t := request.Type; t {
		case "enroll":
			err = r.Enroll(endpoint, c, request.Payload)
			if err != nil {
				_ = c.WriteJSON(NewError(err.Error()))
				_ = c.Close()
				return
			}
		case "action":
			_, err = r.Action(endpoint, request.Payload)
			if err != nil {
				_ = c.WriteJSON(NewError(err.Error()))
				_ = c.Close()
				return
			}
		}
	}
}

func (r *Runtime) findEndpoint(id string) (Endpoint, error) {
	endpoint := Endpoint{}
	err := r.server.Database().Model(&Endpoint{}).Where("id = ?", id).First(&endpoint).Error
	if err != nil {
		return Endpoint{}, err
	}
	return endpoint, nil
}

// socketCloseHandler is called when a websocket connection is terminated
func (r *Runtime) socketCloseHandler(endpoint Endpoint) func(code int, text string) error {
	return func(code int, text string) error {
		r.Unenroll(endpoint)
		return nil
	}
}

// Update sends data to all enrolled endpoints
func (r *Runtime) Update() {
	for id, session := range r.endpoints {
		r.UpdateOne(id, session)
	}
}

func (r *Runtime) UpdateOne(id string, session *Session) {
	db := r.server.Database()

	var model Endpoint
	err := db.Model(&Endpoint{}).Preload("Instances.Module").Where("id = ?", id).First(&model).Error
	if err != nil {
		logger.Error(err.Error())
	}

	response := struct {
		Data map[string]interface{} `json:"data"`
	}{
		Data: map[string]interface{}{},
	}

	for _, instance := range session.Instances {
		i := Instance{}
		err := db.Model(&Instance{}).Preload("Module").Where("id = ?", instance).First(&i).Error
		if err != nil {
			break
		}
		poll, err := i.Poll()
		if err != nil {
			fmt.Println(err)
			break
		}
		response.Data[instance] = poll
	}

	marshal, err := json.Marshal(response)
	if err != nil {
		return
	}

	err = session.Connection.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func (r *Runtime) Action(endpoint Endpoint, payload json.RawMessage) (string, error) {
	action := ActionPayload{}
	err := json.Unmarshal(payload, &action)
	if err != nil {
		return "", err
	}
	instance, err := endpoint.GetInstance(action.Instance, r.server.Database())
	if err != nil {
		return "", err
	}
	run, err := instance.Run(action.Action)
	if err != nil {
		return "", err

	}
	return run, nil
}
