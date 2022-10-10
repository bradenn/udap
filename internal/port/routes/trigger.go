// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/core/domain"
	"udap/internal/core/ports"
)

type triggerRouter struct {
	service ports.TriggerService
}

func (r *triggerRouter) RouteInternal(router chi.Router) {
	router.Post("/triggers/create", r.create)
}

func (r *triggerRouter) RouteExternal(_ chi.Router) {

}

func NewTriggerRouter(service ports.TriggerService) Routable {
	return &triggerRouter{
		service: service,
	}
}

func (r *triggerRouter) create(w http.ResponseWriter, req *http.Request) {

	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(req.Body)
	defer req.Body.Close()
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	sr := domain.Trigger{}

	err = json.Unmarshal(buf.Bytes(), &sr)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	err = r.service.Create(&sr)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	_, err = w.Write([]byte("OK"))
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(200)
}
