// Copyright (c) 2022 Braden Nicholson

package device

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"udap/internal/core/domain"
	"udap/internal/core/ports"
	"udap/internal/port/routes"
)

type router struct {
	service ports.DeviceService
}

func newRouter(service ports.DeviceService) routes.Routable {
	return router{
		service: service,
	}
}

func (r router) RouteInternal(router chi.Router) {
	router.Route("/devices", func(local chi.Router) {
		local.Post("/update", r.update)
	})
}

func (r router) RouteExternal(_ chi.Router) {

}

func (r router) update(w http.ResponseWriter, req *http.Request) {

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
