// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"udap/internal/core/domain"
	"udap/internal/core/ports"
)

type zoneRouter struct {
	service ports.ZoneService
}

func (r zoneRouter) RouteExternal(_ chi.Router) {

}

func NewZoneRouter(service ports.ZoneService) Routable {
	return zoneRouter{
		service: service,
	}
}

func (r zoneRouter) RouteInternal(router chi.Router) {
	router.Post("/zones/create", r.create)
	router.Route("/zones/{id}", func(local chi.Router) {
		local.Post("/delete", r.delete)
		local.Post("/update", r.modify)
		local.Post("/restore", r.restore)
		local.Post("/pin", r.pin)
		local.Post("/unpin", r.unpin)
		local.Post("/entities/{entityId}/add", r.addEntity)
		local.Post("/entities/{entityId}/remove", r.removeEntity)
	})
}

func (r zoneRouter) addEntity(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "invalid id", 400)
	}

	entity := chi.URLParam(req, "entityId")
	if id == "" {
		http.Error(w, "invalid id", 400)
	}

	err := r.service.AddEntity(id, entity)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not pin zone: %s", err.Error()), 400)
		return
	}

	w.WriteHeader(200)
}

func (r zoneRouter) removeEntity(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "invalid id", 400)
		return
	}

	entity := chi.URLParam(req, "entityId")
	if id == "" {
		http.Error(w, "invalid entity", 400)
		return
	}

	err := r.service.RemoveEntity(id, entity)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not remove zone: %s", err.Error()), 400)
		return
	}

	w.WriteHeader(200)
}

func (r zoneRouter) pin(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "invalid id", 400)
	}

	err := r.service.Pin(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not pin zone: %s", err.Error()), 400)
		return
	}

	w.WriteHeader(200)
}

func (r zoneRouter) unpin(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "invalid id", 400)
	}

	err := r.service.Unpin(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not unpin zone: %s", err.Error()), 400)
		return
	}

	w.WriteHeader(200)
}

func (r zoneRouter) delete(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "invalid id", 400)
		return
	}

	err := r.service.Delete(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not delete zone: %s", err.Error()), 400)
		return
	}

	w.WriteHeader(200)
}

func (r zoneRouter) modify(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "invalid id", 400)
	}

	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(w, "could not parse zone", 400)
		return
	}

	ref := domain.Zone{}
	err = json.Unmarshal(buf.Bytes(), &ref)
	if err != nil {
		http.Error(w, "could not parse zone", 400)
		return
	}

	err = r.service.Update(&ref)
	if err != nil {
		http.Error(w, "zone creation failed", 400)
	}

	w.WriteHeader(200)
}

func (r zoneRouter) restore(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "invalid id", 400)
	}

	err := r.service.Restore(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not restore zone: %s", err.Error()), 400)
		return
	}

	w.WriteHeader(200)
}

func (r zoneRouter) create(w http.ResponseWriter, req *http.Request) {

	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(w, "could not parse zone", 400)
		return
	}

	ref := domain.Zone{}
	err = json.Unmarshal(buf.Bytes(), &ref)
	if err != nil {
		http.Error(w, "could not parse zone", 400)
		return
	}

	err = r.service.Create(&ref)
	if err != nil {
		http.Error(w, "zone creation failed", 400)
	}

	w.WriteHeader(200)
}
