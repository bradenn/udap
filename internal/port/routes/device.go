// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/core/domain"
)

type DeviceRouter interface {
	RouteDevices(router chi.Router)
}

type deviceRouter struct {
	service domain.DeviceService
}

func NewDeviceRouter(service domain.DeviceService) DeviceRouter {
	return deviceRouter{
		service: service,
	}
}

func (r deviceRouter) RouteDevices(router chi.Router) {
	router.Route("/devices", func(local chi.Router) {
		local.Post("/update", r.update)
	})
}

func (r deviceRouter) update(w http.ResponseWriter, req *http.Request) {

	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(w, "could not parse device", 400)
		return
	}

	ref := domain.Device{}
	err = json.Unmarshal(buf.Bytes(), &ref)
	if err != nil {
		http.Error(w, "could not parse device", 400)
		return
	}

	err = r.service.Update(&ref)
	if err != nil {
		http.Error(w, "could not find device", 400)
		return
	}

	w.WriteHeader(200)
}
