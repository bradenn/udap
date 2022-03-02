// Copyright (c) 2022 Braden Nicholson

package endpoint

import (
	"github.com/go-chi/chi"
	"net/http"
)

func Router(router chi.Router) {
	router.Post("/", httpPost)
	router.Put("/", httpPut)
	router.Delete("/", httpDelete)
}

func httpPost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func httpDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func httpPut(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
