// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/core/ports"
	"udap/internal/log"
)

type attributeRouter struct {
	service ports.AttributeService
}

func (r *attributeRouter) RouteInternal(router chi.Router) {
	router.Post("/entities/{id}/attributes/{key}/request", r.request)
}

func (r *attributeRouter) RouteExternal(_ chi.Router) {

}

func NewAttributeRouter(service ports.AttributeService) Routable {
	return &attributeRouter{
		service: service,
	}
}

func (r *attributeRouter) request(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	key := chi.URLParam(req, "key")
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(req.Body)
	defer req.Body.Close()
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
	}
	if id != "" && key != "" {
		err = r.service.Request(id, key, buf.String())
		if err != nil {
			log.ErrF(err, "Funny Business:")
		}
	}
	w.Write([]byte("OK"))
	w.WriteHeader(200)
}
