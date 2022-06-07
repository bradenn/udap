// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/core/domain"
)

type ZoneRouter interface {
	RouteZones(router chi.Router)
}

type zoneRouter struct {
	service domain.ZoneService
}

func NewZoneRouter(service domain.ZoneService) ZoneRouter {
	return zoneRouter{
		service: service,
	}
}

func (r zoneRouter) RouteZones(router chi.Router) {
	router.Post("/zones/create", r.create)
	router.Route("/zones/{id}", func(local chi.Router) {
		local.Post("/delete", r.create)
	})
}

func (r zoneRouter) delete(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "invalid id", 400)
	}

	zone, err := r.service.FindById(id)
	if err != nil {
		http.Error(w, "could not parse zone", 400)
		return
	}

	err = r.service.Delete(zone)
	if err != nil {
		http.Error(w, "could not delete zone", 400)
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
