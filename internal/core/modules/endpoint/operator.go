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
	if c.Active() && c.edit != nil {
		c.edit <- body
	}
}

func NewConnection(ws *websocket.Conn) *Connection {
	ch := make(chan any, 8)
	d := make(chan bool)
	a := true
	c := &Connection{
		WS:     ws,
		edit:   ch,
		done:   d,
		active: &a,
	}
	return c
}

func (c *Connection) Close() {
	if c.edit != nil {
		return
	}
	err := c.WS.Close()
	if err != nil {
		return
	}
	close(c.edit)
	close(c.done)
	a := false
	c.active = &a
}

func (c *Connection) Watch() {
	for a := range c.edit {
		if c.WS == nil {
			return
		}
		err := c.WS.WriteJSON(a)
		if err != nil {
			log.Err(err)
			continue
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
	connection.Send(Response{
		Status:    "success",
		Operation: operation,
		Body:      payload,
	})
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

	connection := NewConnection(conn)
	err := m.setConnection(endpoint.Id, connection)
	if err != nil {
		return err
	}

	info, err := systemInfo()
	if err != nil {
		return err
	}

	connection.Send(Response{
		Id:        "",
		Status:    "success",
		Operation: "metadata",
		Body:      Metadata{System: info},
	})
	if err != nil {
		return err
	}

	log.Event("Endpoint '%s' connected.", endpoint.Name)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		connection.Watch()
	}()
	err = m.controller.EmitAll()
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

func (m *endpointOperator) Unenroll(id string) error {
	connection, err := m.getConnection(id)
	if err != nil {
		return err
	}
	connection.Close()
	log.Event("Endpoint '%s' disconnected.", id)
	err = m.removeConnection(id)
	if err != nil {
		return err
	}
	return nil
}

func NewOperator(controller *controller.Controller) domain.EndpointOperator {
	return &endpointOperator{
		controller:  controller,
		connections: map[string]*Connection{},
	}
}
