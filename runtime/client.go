package runtime

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"udap/config"
	"udap/types"
)

// Enroll adds an endpoint's websocket connection to the active update queue
func (r *Runtime) Enroll(endpoint types.Endpoint, conn *websocket.Conn, payload json.RawMessage) error {
	// Initialize a structure to contain the Enrollment instances
	data := EnrollmentPayload{}
	// Attempt to parse the payload into the structure
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return errors.New("invalid enrollment payload")
	}
	// If the endpoint is already enrolled, unenroll
	if r.enrolled(endpoint) {
		r.Unenroll(endpoint)
	}
	// Run the socketCloseHandler when the endpoint is disconnected
	conn.SetCloseHandler(r.socketCloseHandler(endpoint))
	// Generate a session struct, containing the websocket connection, and the requested watchlist
	endpoint.Connection = conn
	endpoint.Subscriptions = data.Instances
	// Enroll the endpoint via the private function
	err = r.enroll(endpoint)
	if err != nil {
		return err
	}
	// Log the new state
	config.Info("Enrolled endpoint '%s'", endpoint.Name)
	// Return a clean bill of health
	err = r.Metadata(endpoint)
	if err != nil {
		return err
	}
	return nil
}

func (r *Runtime) enroll(endpoint types.Endpoint) (err error) {
	// Enroll the endpoint into the runtime
	r.endpoints[endpoint.Id.String()] = &endpoint

	for _, in := range endpoint.Subscriptions {
		if r.instances[in] == nil {
			r.instances[in] = &types.Instance{}
			err = r.instances[in].Instantiate(in)
			if err != nil {
				return
			}
			err = r.instances[in].Load(types.NewAgent(in, r.updater))
			if err != nil {
				return
			}
		}
		r.instances[in].Subscribe(endpoint.Id.String())
	}

	// push all clients to include the new endpoint
	r.push()
	return nil
}

// enrolled returns true if the provided endpoint's UUID is in enrolled
func (r *Runtime) enrolled(endpoint types.Endpoint) bool {
	return r.endpoints[endpoint.Id.String()] != nil
}

// Unenroll removes an endpoint from the active update queue
func (r *Runtime) Unenroll(endpoint types.Endpoint) {
	session := r.endpoints[endpoint.Id.String()]
	// Handle the base case
	if session == nil {
		return
	}
	for _, in := range endpoint.Subscriptions {
		r.instances[in].Unsubscribe(endpoint.Id.String())
	}
	for in, instance := range r.instances {
		if len(instance.Subscribers()) == 0 {
			delete(r.instances, in)
		}
	}
	// Remove the connection reference from the endpoint map
	delete(r.endpoints, endpoint.Id.String())
	// Log the deactivation
	config.Info("Deactivated endpoint '%s'", endpoint.Name)
}

func (r *Runtime) Metadata(endpoint types.Endpoint) error {
	response := Response{
		Status:    "success",
		Operation: METADATA,
		Body:      map[string]interface{}{},
	}

	var modules []types.Module
	findEndpoint, err := r.findEndpoint(endpoint.Id.String())
	if err != nil {
		return err
	}
	var entities []types.Entity
	for _, instance := range endpoint.Instances {
		ent, err := instance.Entities()
		if err != nil {
			return err
		}
		entities = append(entities, ent...)
	}

	response.Body["endpoint"] = findEndpoint
	response.Body["entities"] = entities
	response.Body["modules"] = modules

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
