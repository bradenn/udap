// Copyright (c) 2021 Braden Nicholson

package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var Upgrader websocket.Upgrader

func init() {
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}
