// Copyright (c) 2021 Braden Nicholson

package endpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"udap/module"
	"udap/server"
	"udap/types"
	"udap/udap"
	"udap/udap/db"
	"udap/udap/store"
	"udap/ws"
)

type Request struct {
	Target    string          `json:"target"`
	Operation string          `json:"operation"`
	Body      json.RawMessage `json:"body"`
	Sender    Endpoint
}

type Response struct {
	Status    string                 `json:"status"`
	Operation string                 `json:"operation"`
	Body      map[string]interface{} `json:"body"`
}

// Load modules and instances

// Enroll adds an endpoint's websocket connection to the active update queue
func (e *Endpoint) Enroll(endpoint Endpoint, conn *websocket.Conn, payload json.RawMessage) error {
	// Initialize a structure to contain the Enrollment instances
	data := EnrollmentPayload{}
	// Attempt to parse the payload into the structure
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return errors.New("invalid enrollment payload")
	}
	// If the endpoint is already enrolled, unenroll
	if e.enrolled {
		e.Unenroll()
	}
	// Run the socketCloseHandler when the endpoint is disconnected
	conn.SetCloseHandler(e.socketCloseHandler(endpoint))
	// Generate a session struct, containing the websocket connection, and the requested watchlist

	e.connection = conn
	// Log the new state
	udap.Log("Endpoint '%s' loaded.", endpoint.Name)
	// Return a clean bill of health
	err = e.Metadata()
	if err != nil {
		return err
	}
	return nil
}

// Unenroll removes an endpoint from the active update queue
func (e *Endpoint) Unenroll() {
	// DOne
	udap.Info("Endpoint '%s' unloaded.", e.Name)
}

func writeError(c *websocket.Conn, operation string, err error) {
	response := Response{
		Status:    ERROR,
		Operation: operation,
		Body:      map[string]interface{}{"message": err.Error()},
	}
	_ = c.WriteJSON(response)
	_ = c.Close()
}

const (
	ENROLL    = "enroll"
	METADATA  = "metadata"
	POLL      = "poll"
	RUN       = "run"
	CREATE    = "create"
	MODIFY    = "modify"
	INSTANCE  = "instance"
	ENDPOINT  = "endpoint"
	SUBSCRIBE = "subscribe"
	RESET     = "reset"
	ERROR     = "error"
	SUCCESS   = "success"
	UPDATE    = "update"
	DELETE    = "delete"
)

type EnrollmentPayload struct {
	Instances []string `json:"instances"`
}

type ActionPayload struct {
	Instance string `json:"instance"`
	Action   string `json:"action"`
}

type ModifyPayload struct {
	Instance    string `json:"instance"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ResetPayload struct {
	Instance string `json:"instance"`
}

// RequestHandler upgrades a http request and initiates a websocket session
func HandleSockets(w http.ResponseWriter, req *http.Request) {
	var err error
	// Convert the basic GET request into a WebSocket session
	c, err := ws.Upgrader.Upgrade(w, req, nil)
	tokenParam := chi.URLParam(req, "token")
	// Defer the termination of the session to function return
	defer func(c *websocket.Conn) {
		_ = c.Close()
	}(c)
	// Ignore malformed requests
	if err != nil {
		udap.Err(err)
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
		id, err := server.AuthToken(tokenParam)
		// Aggressively notify the client of their critical mistake
		if err != nil {
			writeError(c, "authentication", err)
			break
		}
		// Attempt to find the endpoint based on its UUID
		endpoint, err := findEndpoint(id)
		if err != nil {
			writeError(c, "authentication", err)
			break
		}
		// Determine the request's intention before assigning its proper function
		endpoint.connection = c
		switch t := request.Target; t {
		case ENDPOINT:
			err = endpoint.EndpointRequest(request, c)
			if err != nil {
				writeError(c, request.Operation, err)
			}
		case INSTANCE:
			err = endpoint.InstanceRequest(request)
			if err != nil {
				writeError(c, request.Operation, err)
			}
		}
		err = endpoint.Metadata()
		if err != nil {
			writeError(c, request.Operation, err)
		}
	}
}

func (e *Endpoint) Metadata() error {
	response := Response{
		Status:    "success",
		Operation: METADATA,
		Body:      map[string]interface{}{},
	}

	var modules []module.Module
	ep, err := findEndpoint(e.Id)
	if err != nil {
		return err
	}
	var entities []types.Entity
	for _, instance := range e.Instances {
		ent, err := instance.Entities()
		if err != nil {
			return err
		}
		for i, entity := range ent {
			get, err := store.Get(fmt.Sprintf("instance.%s.entity.%s.state", entity.InstanceId.String(), entity.Name))
			if err != nil {
				continue
			} else {
				entity.State = get
				ent[i] = entity
			}

		}
		entities = append(entities, ent...)
	}

	response.Body["endpoint"] = ep
	response.Body["entities"] = entities
	response.Body["modules"] = modules

	marshal, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = e.connection.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		log.Println("write:", err)
		return err
	}
	return nil
}

// push sends data to all enrolled endpoints
func (e *Endpoint) Push(instances map[string]interface{}) {

	response := Response{
		Status:    SUCCESS,
		Operation: UPDATE,
		Body:      map[string]interface{}{},
	}

	for _, instance := range e.Subscriptions {
		response.Body[instance.UUID()] = instances[instance.UUID()]
	}

	marshal, err := json.Marshal(response)
	if err != nil {
		return
	}

	err = e.connection.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		log.Println("write:", err)
		return
	}

}

func findEndpoint(id string) (Endpoint, error) {
	ep := Endpoint{}
	err := db.DB.Model(&Endpoint{}).Preload("Instances.Module").Where("id = ?",
		id).First(&ep).Error
	if err != nil {
		return ep, err
	}
	return ep, nil
}

// socketCloseHandler is called when a websocket connection is terminated
func (e *Endpoint) socketCloseHandler(endpoint Endpoint) func(code int, text string) error {
	return func(code int, text string) error {

		return nil
	}
}
