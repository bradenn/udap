// Copyright (c) 2022 Braden Nicholson

package endpoint

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/core/ports"
	"udap/internal/log"
)

const (
	ENROLL   = "enroll"
	UNENROLL = "unenroll"
)

type endpointOperation struct {
	operation string
	endpoint  *domain.Endpoint
	response  chan error
}

func (e *endpointOperation) Respond(err error) {
	select {
	case e.response <- err:
	default:
	}
}

func newOperation(operation string, endpoint *domain.Endpoint) (endpointOperation, chan error) {
	response := make(chan error)
	return endpointOperation{
		response:  response,
		operation: operation,
		endpoint:  endpoint,
	}, response
}

type endpointOperator struct {
	local         map[string]*domain.Endpoint
	localChannel  chan endpointOperation
	localTransmit chan Response
	done          chan bool
	controller    *controller.Controller
}

func NewOperator(controller *controller.Controller) ports.EndpointOperator {

	op := &endpointOperator{
		controller:    controller,
		local:         map[string]*domain.Endpoint{},
		done:          make(chan bool),
		localChannel:  make(chan endpointOperation),
		localTransmit: make(chan Response, 32),
	}

	go func() {
		err := op.listen()
		if err != nil {
			log.Err(err)
		}
	}()

	return op
}

func (m *endpointOperator) handleShutdown() {
	records := m.local
	for _, endpoint := range records {
		_ = endpoint.Connection.WriteMessage(websocket.CloseGoingAway, nil)
	}
}

func (m *endpointOperator) transmitSingle(transmission Response) error {
	err := m.local[transmission.Endpoint].Connection.WriteJSON(transmission)
	if err != nil {
		return err
	}
	return nil
}

func (m *endpointOperator) handleTransmit(transmission Response) error {
	if transmission.Endpoint != "" {
		return m.transmitSingle(transmission)
	}
	records := m.local
	for _, endpoint := range records {
		err := endpoint.Connection.WriteJSON(transmission)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *endpointOperator) handleOperation(operation endpointOperation) error {
	endpoint := operation.endpoint
	switch operation.operation {
	case ENROLL:
		ref := m.local[endpoint.Id]
		if ref != nil {
			operation.Respond(fmt.Errorf("endpoint already enrolled"))
		}
		m.local[endpoint.Id] = operation.endpoint
		operation.Respond(nil)
		break
	case UNENROLL:
		ref := m.local[endpoint.Id]
		if ref == nil {
			operation.Respond(fmt.Errorf("endpoint is not enrolled"))
		}
		delete(m.local, endpoint.Id)
		operation.Respond(nil)
		break
	default:
		return fmt.Errorf("unknown operation '%s'", operation.operation)
	}
	return nil
}

func (m *endpointOperator) listen() error {
	for {
		select {
		case operation := <-m.localChannel:
			err := m.handleOperation(operation)
			if err != nil {
				log.Err(err)
				continue
			}
		case transmission := <-m.localTransmit:
			err := m.handleTransmit(transmission)
			if err != nil {
				log.Err(err)
				continue
			}
		case <-m.done:
			m.handleShutdown()
			return nil
		}
	}
}

func (m *endpointOperator) enrollEndpoint(endpoint *domain.Endpoint) error {
	operation, errChan := newOperation("enroll", endpoint)
	m.localChannel <- operation
	if err := <-errChan; err != nil {
		return err
	}
	err := m.sendMetadata(endpoint.Id)
	if err != nil {
		return err
	}
	err = m.controller.EmitAll()
	if err != nil {
		log.Err(err)
		return nil
	}
	log.Event("Endpoint '%s' enrolled.", endpoint.Name)
	return nil
}

func (m *endpointOperator) unenrollEndpoint(endpoint *domain.Endpoint) error {
	operation, errChan := newOperation("unenroll", endpoint)
	m.localChannel <- operation
	if err := <-errChan; err != nil {
		return err
	}
	log.Event("Endpoint '%s' unenrolled.", endpoint.Name)
	return nil
}

func (m *endpointOperator) CloseAll() error {
	m.done <- true
	return nil
}

func (m *endpointOperator) SendAll(id string, operation string, payload any) error {

	transmission := Response{
		Endpoint:  "",
		Id:        id,
		Status:    "success",
		Operation: operation,
		Body:      payload,
	}
	timer := time.NewTimer(time.Millisecond * 500)
	select {
	// Attempt to push the payload to the channel
	case m.localTransmit <- transmission:
		// Cancel the timer if payload is sent
		timer.Stop()
		// Exit normally
		return nil
	case <-timer.C:
		log.Event("transit transmission timed out")
		// Exit quietly if the payload could not be sent
		return nil
	}

}

func (m *endpointOperator) Send(id string, operation string, payload any) error {

	transmission := Response{
		Endpoint:  id,
		Id:        id,
		Status:    "success",
		Operation: operation,
		Body:      payload,
	}

	m.localTransmit <- transmission

	return nil
}

type Metadata struct {
	System System `json:"system"`
}

type Response struct {
	Endpoint  string `json:"endpoint"`
	Id        string `json:"id"`
	Status    string `json:"status"`
	Operation string `json:"operation"`
	Body      any    `json:"body"`
}

func (m *endpointOperator) sendMetadata(id string) error {
	info, err := systemInfo()
	if err != nil {
		return err
	}
	m.localTransmit <- Response{
		Endpoint:  id,
		Id:        "",
		Status:    "success",
		Operation: "metadata",
		Body:      Metadata{System: info},
	}
	return nil
}

func (m *endpointOperator) Enroll(endpoint *domain.Endpoint, conn *websocket.Conn) error {
	if endpoint.Connection != nil {
		return fmt.Errorf("connection already exists")
	}
	endpoint.Connection = conn
	endpoint.Connected = true
	err := m.enrollEndpoint(endpoint)
	if err != nil {
		return err
	}
	return nil

}

func (m *endpointOperator) Unenroll(endpoint *domain.Endpoint) error {
	err := m.unenrollEndpoint(endpoint)
	if err != nil {
		return err
	}
	endpoint.Connection = nil
	endpoint.Connected = false
	return nil
}
