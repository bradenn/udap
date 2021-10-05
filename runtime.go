package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"udap/logger"
	"udap/server"
)

const (
	pathFmt   = "./plugins"
	pluginFmt = pathFmt + "/%s/%s.so"
)

const (
	ERROR    = "error"
	ENROLL   = "enroll"
	METADATA = "metadata"
	POLL     = "poll"
	ACTION   = "action"
)

type Session struct {
	Connection *websocket.Conn
	Instances  []string
}

type Response struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func ErrorResponse(message string) string {
	r := Response{
		Type:    ERROR,
		Payload: map[string]interface{}{},
	}
	r.Payload["message"] = message
	marshal, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(marshal)
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

// Runtime manages the event-loop for all instances, as well as the websocket connections between UDAP and endpoints.
type Runtime struct {
	endpoints map[string]*Session
	instances map[string]*Instance
	// subscribers contains an array of endpoint ids for each instance id
	subscribers map[string][]string
	// cache holds the current data for each live instance
	cache    map[string]string
	server   server.Server
	upgrader websocket.Upgrader
}

// NewRuntime creates, initializes, and configures a websocket runtime
func NewRuntime(server server.Server) Runtime {
	runtime := Runtime{
		endpoints:   map[string]*Session{},
		instances:   map[string]*Instance{},
		subscribers: map[string][]string{},
		cache:       map[string]string{},
		server:      server,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	server.Router().HandleFunc("/ws", runtime.EndpointHandler)
	return runtime
}

func (r *Runtime) Begin(duration time.Duration) {
	go func() {
		for {
			r.tick()
			time.Sleep(duration)
		}
	}()
}

func (r *Runtime) tick() {
	for instanceId, endpointArray := range r.subscribers {
		if len(endpointArray) >= 1 {
			r.updateInstance(instanceId)
		}
	}
	r.push()
}

// discoverModules locates and initializes modules in the plugin directory
func (r *Runtime) discoverModules() {
	dir, err := ioutil.ReadDir(pathFmt)
	if err != nil {
		return
	}

	for _, info := range dir {
		if info.IsDir() {
			mod := Module{}
			mod.Path = info.Name()
			err = mod.Load(r.server.Database())
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}
}

// Enroll adds an endpoint's websocket connection to the active update queue
func (r *Runtime) Enroll(endpoint Endpoint, conn *websocket.Conn, payload json.RawMessage) error {
	if r.enrolled(endpoint) {
		return errors.New("endpoint is already enrolled")
	}
	// Initialize a structure to contain the Enrollment instances
	data := EnrollmentPayload{}
	// Attempt to parse the payload into the structure
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return errors.New("invalid enrollment payload")
	}
	// Log the new state
	logger.Info("Enrolled endpoint '%s'", endpoint.Name)
	// Set a close handler to deactivate the endpoint when the connection is severed
	conn.SetCloseHandler(r.socketCloseHandler(endpoint))
	// Generate a session struct, containing the websocket connection, and the requested watchlist
	session := Session{
		Connection: conn,
		Instances:  data.Instances,
	}
	for _, instance := range data.Instances {
		r.subscribe(endpoint.Id.String(), instance)
	}
	// Enroll the endpoint via the private function
	r.enroll(endpoint.Id.String(), session)
	// Return a clean bill of health
	err = r.Metadata(endpoint)
	if err != nil {
		return err
	}
	return nil
}

func (r *Runtime) subscribe(endpointId string, instanceId string) {
	subs := r.subscribers[instanceId]
	if subs == nil {
		subs = []string{}
	}
	for _, sub := range subs {
		if sub == endpointId {
			return
		}
	}
	subs = append(subs, endpointId)
	r.subscribers[instanceId] = subs
}

func (r *Runtime) unsubscribe(endpointId string, instanceId string) {
	subs := r.subscribers[instanceId]
	if subs == nil || len(subs) == 0 {
		return
	}
	for index, sub := range subs {
		if sub == endpointId {
			subs = append(subs[:index], subs[index+1:]...)
			return
		}
	}
	r.subscribers[instanceId] = subs
}

func (r *Runtime) enroll(endpointId string, session Session) {
	// Enroll the endpoint into the runtime
	r.endpoints[endpointId] = &session
	// push all clients to include the new endpoint
	r.push()
}

// enrolled returns true if the provided endpoint's UUID is in enrolled
func (r *Runtime) enrolled(endpoint Endpoint) bool {
	return r.endpoints[endpoint.Id.String()] != nil
}

// Unenroll removes an endpoint from the active update queue
func (r *Runtime) Unenroll(endpoint Endpoint) {
	session := r.endpoints[endpoint.Id.String()]
	// Handle the base case
	if session == nil {
		return
	}

	for _, instance := range session.Instances {
		r.unsubscribe(endpoint.Id.String(), instance)
	}
	// Remove the connection reference from the endpoint map
	delete(r.endpoints, endpoint.Id.String())
	// Log the deactivation
	logger.Info("Deactivated endpoint '%s'", endpoint.Name)
}

// EndpointHandler upgrades a http request and initiates a websocket session
func (r *Runtime) EndpointHandler(w http.ResponseWriter, req *http.Request) {
	// Convert the basic GET request into a Websocket session
	c, err := r.upgrader.Upgrade(w, req, nil)
	// Defer the termination of the session to function return
	defer func(c *websocket.Conn) {
		_ = c.Close()
	}(c)
	// Ignore malformed requests
	if err != nil {
		logger.Err(err)
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
			logger.Err(err)
			_ = c.WriteJSON(ErrorResponse("Invalid Token"))
			_ = c.Close()
			break
		}
		// Attempt to find the endpoint based on its UUID
		endpoint, err := r.findEndpoint(id)
		if err != nil {
			logger.Err(err)
			_ = c.WriteJSON(ErrorResponse("Invalid Endpoint"))
			_ = c.Close()
			break
		}
		// Determine the request's intention before assigning its proper function
		switch t := request.Type; t {
		case ENROLL:
			err = r.Enroll(endpoint, c, request.Payload)
			if err != nil {
				_ = c.WriteJSON(ErrorResponse(err.Error()))
				_ = c.Close()
				return
			}
		case ACTION:
			_, err = r.Action(endpoint, request.Payload)
			if err != nil {
				_ = c.WriteJSON(ErrorResponse(err.Error()))
				_ = c.Close()
				return
			}
		case METADATA:
			err = r.Metadata(endpoint)
			if err != nil {
				_ = c.WriteJSON(ErrorResponse(err.Error()))
				_ = c.Close()
				return
			}
		}
	}
}

func (r *Runtime) findEndpoint(id string) (Endpoint, error) {
	endpoint := Endpoint{
		Instances: []Instance{},
	}
	err := r.server.Database().Model(&Endpoint{}).Preload("Instances.Module").Where("id = ?", id).First(&endpoint).Error
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

// push sends data to all enrolled endpoints
func (r *Runtime) push() {
	for _, session := range r.endpoints {
		r.pushEndpoint(session)
	}
}

func (r *Runtime) updateInstance(instanceId string) {
	db := r.server.Database()
	i := Instance{}
	err := db.Model(&Instance{}).Preload("Module").Where("id = ?", instanceId).First(&i).Error
	if err != nil {
		logger.Err(err)
		return
	}
	poll, err := i.Poll()
	if err != nil {
		logger.Err(err)
		return
	}
	r.cache[instanceId] = poll
}

func (r *Runtime) pushEndpoint(session *Session) {
	response := Response{
		Type:    POLL,
		Payload: map[string]interface{}{},
	}

	for _, instance := range session.Instances {
		response.Payload[instance] = r.cache[instance]
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

// Action forwards the websocket action message to the defined instance function
func (r *Runtime) Action(endpoint Endpoint, payload json.RawMessage) (string, error) {
	// A variable to hold the payload is initialized
	action := ActionPayload{}
	// Attempt to parse the payload into the ActionPayload struct
	err := json.Unmarshal(payload, &action)
	if err != nil {
		return "", err
	}
	// Find the instance defined in the payload, and authorize the endpoint
	instance, err := endpoint.GetInstance(action.Instance, r.server.Database())
	if err != nil {
		return "", err
	}
	// Run the requested function
	run, err := instance.Run(action.Action)
	if err != nil {
		return "", err

	}
	return run, nil
}

func (r *Runtime) Metadata(endpoint Endpoint) error {
	response := Response{
		Type:    METADATA,
		Payload: map[string]interface{}{},
	}

	response.Payload["endpoint"] = endpoint

	marshal, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = r.endpoints[endpoint.Id.String()].Connection.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		log.Println("write:", err)
		return err
	}
	return nil
}
