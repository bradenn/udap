// Copyright (c) 2021 Braden Nicholson

package types

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Session struct {
	Connection *websocket.Conn
	Endpoint   uuid.UUID
}
