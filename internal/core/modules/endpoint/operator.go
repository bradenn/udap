// Copyright (c) 2022 Braden Nicholson

package endpoint

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/log"
)

type Connection struct {
	WS     *websocket.Conn
	active *bool
	edit   chan any
	done   chan bool
}

func (c *Connection) Active() bool {
	return *c.active
}

func (c *Connection) Send(body any) {
	c.edit <- body
}

func NewConnection(ws *websocket.Conn) *Connection {
	ch := make(chan any, 24)
	d := make(chan bool)
	active := true
	c := &Connection{
		WS:     ws,
		edit:   ch,
		done:   d,
		active: &active,
	}

	ws.SetCloseHandler(func(code int, text string) error {
		c.Close(false)
		return nil
	})

	return c
}

func (c *Connection) Close(self bool) {
	c.done <- self
	active := false
	c.active = &active
}

func (c *Connection) Watch() {
	for {
		select {
		case t := <-c.done:
			if t {
				c.WS.WriteMessage(websocket.CloseGoingAway, nil)
			}
			return
		case req := <-c.edit:
			err := c.WS.WriteJSON(req)
			if err != nil {
				log.Err(err)
				return
			}
		}
	}
}

type endpointOperator struct {
	connections map[string]*Connection
	controller  *controller.Controller
}

func (m *endpointOperator) getConnection(id string) (*Connection, error) {
	ref := m.connections[id]
	if ref == nil {
		return nil, fmt.Errorf("connection not found")
	}
	return ref, nil
}

func (m *endpointOperator) setConnection(id string, connection *Connection) error {
	m.connections[id] = connection
	return nil
}

func (m *endpointOperator) CloseAll() error {
	for _, connection := range m.connections {
		connection.Close(true)
	}
	return nil
}

func (m *endpointOperator) removeConnection(id string) error {
	ref := m.connections[id]
	if ref != nil {
		return nil
	}
	delete(m.connections, id)
	return nil
}

func (m *endpointOperator) SendAll(id string, operation string, payload any) error {
	for _, conn := range m.connections {
		if conn.Active() {
			conn.Send(Response{
				Id:        id,
				Status:    "success",
				Operation: operation,
				Body:      payload,
			})
		}
	}
	return nil
}

func (m *endpointOperator) Send(id string, operation string, payload any) error {
	connection, err := m.getConnection(id)
	if err != nil {
		return err
	}
	if connection.Active() {
		connection.Send(Response{
			Status:    "success",
			Operation: operation,
			Body:      payload,
		})
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

func sendMetadata(connection *Connection) error {
	info, err := systemInfo()
	if err != nil {
		return err
	}
	if connection.Active() {
		connection.Send(Response{
			Id:        "",
			Status:    "success",
			Operation: "metadata",
			Body:      Metadata{System: info},
		})
	}
	return nil
}

func (m *endpointOperator) Enroll(endpoint *domain.Endpoint, conn *websocket.Conn) error {
	// Initialize a new endpoint connection
	connection := NewConnection(conn)
	// Insert the connection into the local map
	err := m.setConnection(endpoint.Id, connection)
	if err != nil {
		return err
	}
	// Send the system metadata to the client
	err = sendMetadata(connection)
	if err != nil {
		return err
	}
	// Create a thread to handle sending data to the endpoint
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		connection.Watch()
	}()
	// Send all the current system context to the endpoint
	err = m.controller.EmitAll()
	if err != nil {
		log.Err(err)
	}
	// Log the current state
	log.Event("Endpoint '%s' connected.", endpoint.Name)
	// Listen for incoming close messages
	for {
		_, _, err = connection.WS.ReadMessage()
		if err != nil {
			break
		}
	}
	// Wait until the watch function exits
	wg.Wait()
	// Remove the connection from the local map
	err = m.removeConnection(endpoint.Id)
	if err != nil {
		return err
	}
	// Print disconnect message
	log.Event("Endpoint '%s' disconnected.", endpoint.Name)
	return nil
}

func (m *endpointOperator) Unenroll(id string) error {

	return nil
}

func NewOperator(controller *controller.Controller) domain.EndpointOperator {
	return &endpointOperator{
		controller:  controller,
		connections: map[string]*Connection{},
	}
}
