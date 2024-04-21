// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"udap/internal/srv/store"
)

type traceRouter struct {
	store *store.Store
}

func (r *traceRouter) RouteInternal(router chi.Router) {
	router.Post("/trace", r.trace)
}

func (r *traceRouter) RouteExternal(_ chi.Router) {

}

func NewTraceRouter(server *store.Store) Routable {
	return &traceRouter{
		store: server,
	}
}

type TraceRequest struct {
	To     int64    `json:"to"`
	From   int64    `json:"from"`
	Window int      `json:"window"`
	Mode   string   `json:"mode"`
	Labels []string `json:"labels"`
}

type Trace struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
	Time   []int64           `json:"time"`
	Data   []float64         `json:"data"`
}

type TraceResults struct {
	Traces []Trace `json:"traces"`
}

func (r *traceRouter) trace(w http.ResponseWriter, req *http.Request) {

	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(req.Body)
	sr := TraceRequest{}
	err = json.Unmarshal(buf.Bytes(), &sr)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	tr := TraceResults{}

	tr.Traces = []Trace{}

	traces, labels, err := r.store.Traces(sr.From, sr.To, sr.Window, sr.Mode, sr.Labels)
	if err != nil {

		return
	}

	for s, m := range traces {
		t := Trace{
			Name:   s,
			Time:   []int64{},
			Data:   []float64{},
			Labels: labels[s],
		}

		for a, b := range m {
			t.Time = append(t.Time, a)
			t.Data = append(t.Data, b)
		}

		tr.Traces = append(tr.Traces, t)
	}

	marshal, err := json.Marshal(tr)
	if err != nil {
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		return
	}

}
