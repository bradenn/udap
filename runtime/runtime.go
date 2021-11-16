package runtime

import (
	"github.com/gorilla/websocket"
	"net/http"
	"udap/server"
	"udap/types"
)

// Runtime manages the event-loop for all instances, as well as the websocket connections between UDAP and endpoints.
type Runtime struct {
	modules map[string]*types.Module
	// endpoints is a map of each active session
	endpoints map[string]*types.Endpoint
	// instances contains all active instances of udap
	instances map[string]*types.Instance
	// cache holds the current data for each live instance
	cache map[string]string

	updater chan types.UpdateBuffer

	server   server.Server
	upgrader websocket.Upgrader
}

// New creates, initializes, and configures a websocket runtime
func New(server server.Server) Runtime {
	// The updater channel is called by modules, we initialize it to queue 16 at a time
	updater := make(chan types.UpdateBuffer, 16)
	runtime := Runtime{
		modules:   map[string]*types.Module{},
		endpoints: map[string]*types.Endpoint{},
		instances: map[string]*types.Instance{},
		cache:     map[string]string{},
		updater:   updater,
		server:    server,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	// Discover any unloaded modules
	runtime.discoverModules()
	// Route WebSocket requests here
	server.Router().HandleFunc("/ws/{token}", runtime.RequestHandler)
	// Set async listener for the updating from modules and daemons
	go runtime.Listen()
	// return the runtime
	return runtime
}
