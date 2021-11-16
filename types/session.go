package types

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Session struct {
	Connection *websocket.Conn
	Endpoint   uuid.UUID
}
