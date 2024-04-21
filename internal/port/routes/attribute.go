// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"udap/internal/core/ports"
)

type attributeRouter struct {
	service ports.AttributeService
}

func (r *attributeRouter) RouteInternal(router chi.Router) {
	router.Post("/entities/{id}/attributes/{key}/request", r.request)
	router.Post("/attribute/{id}/delete", r.delete)
	router.Post("/attribute/summary", r.summary)

}

func (r *attributeRouter) RouteExternal(_ chi.Router) {

}

func NewAttributeRouter(service ports.AttributeService) Routable {
	return &attributeRouter{
		service: service,
	}
}

func (r *attributeRouter) delete(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(req.Body)
	defer req.Body.Close()
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	byId, err := r.service.FindById(id)
	if err != nil {
		return
	}

	if id != "" {
		err = r.service.Delete(byId)
		if err != nil {
			w.Write([]byte(err.Error()))
			//w.WriteHeader(500)
			return
			//log.ErrF(err, "Funny Business:")
		}
	}

	w.Write([]byte("OK"))
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
		return
	}
	if id != "" && key != "" {
		err = r.service.Request(id, key, buf.String())
		if err != nil {
			w.Write([]byte(err.Error()))
			//w.WriteHeader(500)
			return
			//log.ErrF(err, "Funny Business:")
		}
	}
	w.Write([]byte("OK"))
}

type SummaryRequest struct {
	Key    string `json:"key"`
	To     int64  `json:"to"`
	From   int64  `json:"from"`
	Window int    `json:"window"`
	Mode   string `json:"mode"`
}

func (r *attributeRouter) summary(w http.ResponseWriter, req *http.Request) {

	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(req.Body)
	sr := SummaryRequest{}
	err = json.Unmarshal(buf.Bytes(), &sr)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	summary, err := r.service.Summary(sr.Key, sr.From, sr.To, sr.Window, sr.Mode)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	marshal, err := json.Marshal(summary)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	_, _ = w.Write(marshal)
}
