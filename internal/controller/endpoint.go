// Copyright (c) 2022 Braden Nicholson

package controller

import (
	"fmt"
	"github.com/go-chi/chi"
	"sync"
	"udap/internal/bond"
	"udap/internal/log"
	"udap/internal/models"
	"udap/internal/store"
)

type Response struct {
	Id        string      `json:"id"`
	Status    string      `json:"status"`
	Operation string      `json:"operation"`
	Body      interface{} `json:"body"`
}

type Endpoints struct {
	PolyBuffer
	Observable
	bond   *bond.Bond
	router chi.Router
}

func (e *Endpoints) Handle(msg bond.Msg) (res interface{}, err error) {
	switch t := msg.Operation; t {
	case "create":
		return e.create(msg)
	default:
		return nil, fmt.Errorf("operation '%s' is not defined", t)
	}
}

func LoadEndpoints() (m *Endpoints) {
	m = &Endpoints{}
	m.data = sync.Map{}
	m.raw = map[string]interface{}{}
	m.Run()
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

func (e *Endpoints) create(msg bond.Msg) (res interface{}, err error) {
	ep := models.NewEndpoint(msg.Payload)
	err = store.DB.Create(&ep).Error
	if err != nil {
		return nil, err
	}
	e.Set(ep.Id, &ep)
	return ep, err
}

func (e *Endpoints) unenroll(msg bond.Msg) (res interface{}, err error) {
	endpoint := e.Find(msg.Id)
	endpoint.Unenroll()
	return nil, nil
}

func (e *Endpoints) Compile() (endpoints []models.Endpoint, err error) {
	for _, s := range e.Keys() {
		endpoint := e.Find(s)
		endpoints = append(endpoints, *endpoint)
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
	e.emit(id, endpoint)
	e.set(id, endpoint)

}

func (e *Endpoints) Save(endpoint *models.Endpoint) {
	store.DB.Model(&models.Endpoint{}).Save(endpoint)
}
