// Copyright (c) 2022 Braden Nicholson

package endpoint

import (
	"github.com/gorilla/websocket"
	"udap/internal/core/domain"
	"udap/internal/log"
)

type endpointOperator struct {
	connections map[string]*websocket.Conn
}

func (m *endpointOperator) Send(id string, operation string, payload any) error {
	if m.connections[id] == nil {
		return nil
	}
	err := m.connections[id].WriteJSON(Response{
		Status:    "success",
		Operation: operation,
		Body:      payload,
	})
	if err != nil {
		return err
	}
	return nil
}

type Metadata struct {
	System System `json:"system"`
}

type Response struct {
	Id        string `json:"id"`
	Status    string `json:"status"`
	Operation string `json:"operation"`
	Body      any    `json:"body"`
}

func (m *endpointOperator) Enroll(endpoint *domain.Endpoint, conn *websocket.Conn) error {
	m.connections[endpoint.Id] = conn
	m.connections[endpoint.Id].SetCloseHandler(func(code int, text string) error {
		log.Event("Endpoint '%s' disconnected.", endpoint.Name)
		return nil
	})
	info, err := systemInfo()
	if err != nil {
		return err
	}
	err = conn.WriteJSON(Response{
		Id:        "",
		Status:    "success",
		Operation: "metadata",
		Body:      Metadata{System: info},
	})
	if err != nil {
		return err
	}

	log.Event("Endpoint '%s' connected.", endpoint.Name)
	return nil
}

func NewOperator() domain.EndpointOperator {
	return &endpointOperator{
		connections: map[string]*websocket.Conn{},
	}
}
