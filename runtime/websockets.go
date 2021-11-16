package runtime

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"net/http"
	"udap/config"
	"udap/server"
	"udap/types"
)

const (
	ENROLL   = "enroll"
	METADATA = "metadata"
	POLL     = "poll"
	RUN      = "run"
	CREATE   = "create"
	MODIFY   = "modify"
	INSTANCE = "instance"
	ENDPOINT = "endpoint"
	RESET    = "reset"
	ERROR    = "error"
	SUCCESS  = "success"
	UPDATE   = "update"
	DELETE   = "delete"
)

type Request struct {
	Target    string          `json:"target"`
	Operation string          `json:"operation"`
	Body      json.RawMessage `json:"body"`
	Sender    types.Endpoint
}

type Response struct {
	Status    string                 `json:"status"`
	Operation string                 `json:"operation"`
	Body      map[string]interface{} `json:"body"`
}

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
func (r *Runtime) RequestHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	// Convert the basic GET request into a WebSocket session
	c, err := r.upgrader.Upgrade(w, req, nil)
	tokenParam := chi.URLParam(req, "token")
	// Defer the termination of the session to function return
	defer func(c *websocket.Conn) {
		_ = c.Close()
	}(c)
	// Ignore malformed requests
	if err != nil {
		config.Err(err)
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
		endpoint, err := r.findEndpoint(id)
		if err != nil {
			writeError(c, "authentication", err)
			break
		}
		request.Sender = endpoint
		// Determine the request's intention before assigning its proper function

		switch t := request.Target; t {
		case ENDPOINT:
			err = r.EndpointRequest(request, c)
			if err != nil {
				writeError(c, request.Operation, err)
			}
		case INSTANCE:
			err = r.InstanceRequest(request)
			if err != nil {
				writeError(c, request.Operation, err)
			}
		}
		err = r.Metadata(request.Sender)
		if err != nil {
			writeError(c, request.Operation, err)
		}
	}
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

func (r *Runtime) EndpointRequest(request Request, c *websocket.Conn) (err error) {
	switch t := request.Operation; t {
	case METADATA:
		err = r.Metadata(request.Sender)
		if err != nil {
			return err
		}
	case ENROLL:
		err = r.Enroll(request.Sender, c, request.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

type IdentifierBody struct {
	Id string `json:"id"`
}

func (r *Runtime) InstanceRequest(request Request) (err error) {
	body := IdentifierBody{}
	instance := types.Instance{}

	err = json.Unmarshal(request.Body, &body)
	if err != nil {
		return err
	}

	if body.Id == "" {
		switch t := request.Operation; t {
		case CREATE:
			err = instance.Create(request.Body)
			if err != nil {
				return err
			}
			err = request.Sender.GrantInstance(instance)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid operation '%s'; you may have forgotten to provice an Id", t)
		}
		return nil
	}

	instance, err = request.Sender.GetInstance(body.Id)
	if instance.ModuleId.String() == "" {
		panic("Wtf? No Module ID on this bad boy")
	}
	instance.Module = r.modules[instance.ModuleId.String()]
	switch t := request.Operation; t {
	case MODIFY:
		err = instance.Modify(request.Body)
		if err != nil {
			return err
		}
	case RUN:
		err = instance.Run(string(request.Body))
		if err != nil {
			return err
		}
	case RESET:
		err = instance.Reset()
		if err != nil {
			return err
		}
	case DELETE:
		err = instance.Delete(request.Sender)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid operation '%s'", t)
	}

	return nil
}

func (r *Runtime) findEndpoint(id string) (types.Endpoint, error) {
	endpoint := types.Endpoint{
		Instances: []types.Instance{},
	}
	err := r.server.Database.Model(&types.Endpoint{}).Preload("Instances.Module").Where("id = ?", id).First(&endpoint).Error
	if err != nil {
		return types.Endpoint{}, err
	}
	return endpoint, nil
}

// socketCloseHandler is called when a websocket connection is terminated
func (r *Runtime) socketCloseHandler(endpoint types.Endpoint) func(code int, text string) error {
	return func(code int, text string) error {
		r.Unenroll(endpoint)
		return nil
	}
}
