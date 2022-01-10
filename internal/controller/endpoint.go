// Copyright (c) 2022 Braden Nicholson

package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"sync"
	"udap/internal/bond"
	"udap/internal/log"
	"udap/internal/models"
	"udap/internal/store"
)

type Request struct {
	Target    string          `json:"target"`
	Operation string          `json:"operation"`
	Body      json.RawMessage `json:"body"`
	Sender    string
}

type Response struct {
	Status    string `json:"status"`
	Operation string `json:"operation"`
	Body      any    `json:"body"`
}

type Endpoints struct {
	PolyBuffer
	bond   *bond.Bond
	router chi.Router
}

func (e *Endpoints) Handle(msg bond.Msg) (res any, err error) {
	switch t := msg.Operation; t {
	case "compile":
		return e.compile(msg)
	case "enroll":
		return e.enroll(msg)
	case "unenroll":
		return e.enroll(msg)
	default:
		return nil, fmt.Errorf("operation '%s' is not defined", t)
	}
}

func LoadEndpoints() (m *Endpoints) {
	m = &Endpoints{}
	m.data = sync.Map{}
	m.raw = map[string]any{}
	m.FetchAll()
	return m
}

func (e *Endpoints) FetchAll() {
	var endpoints []*models.Endpoint
	store.DB.Model(&models.Endpoint{}).Find(&endpoints)
	for _, endpoint := range endpoints {
		e.Set(endpoint.Id, endpoint)
	}
}

func (e *Endpoints) unenroll(msg bond.Msg) (res any, err error) {
	endpoint := e.Find(msg.Id)
	endpoint.Unenroll()
	return nil, nil
}

func (e *Endpoints) enroll(msg bond.Msg) (res any, err error) {
	endpoint := e.Find(msg.Id)
	err = endpoint.Enroll()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Endpoints) compile(msg bond.Msg) (res any, err error) {
	var endpoints []models.Endpoint
	for _, s := range e.Keys() {
		endpoint := e.Find(s)
		endpoints = append(endpoints, *endpoint)
	}
	return endpoints, nil
}

func (e *Endpoints) Compile() (endpoints []map[string]any, err error) {
	for _, s := range e.Keys() {
		endpoint := e.Find(s)
		marshal, err := json.Marshal(endpoint)
		if err != nil {
			return nil, err
		}
		var cache map[string]any
		err = json.Unmarshal(marshal, &cache)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, cache)
	}
	return endpoints, err
}

func (e *Endpoints) Find(id string) *models.Endpoint {

	dat := e.get(id)
	if dat == struct{}{} {
		dat = nil
	}
	if dat == nil {
		endpoint := &models.Endpoint{}
		endpoint.Id = id
		err := endpoint.Fetch()
		if err != nil {
			log.Err(err)
		}
		return endpoint
	}

	return dat.(*models.Endpoint)
}

func (e *Endpoints) Set(id string, endpoint *models.Endpoint) {

	e.set(id, endpoint)

}

func (e *Endpoints) Save(endpoint *models.Endpoint) {
	store.DB.Model(&models.Endpoint{}).Save(endpoint)
}
