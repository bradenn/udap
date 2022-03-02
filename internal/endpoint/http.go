// Copyright (c) 2022 Braden Nicholson

package endpoint

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

// HttpRouter distributes rest requests to the appropriate functions
func HttpRouter(router chi.Router) {
	router.Post("/", httpPost)
	router.Put("/", httpPut)
	router.Delete("/", httpDelete)
}

// httpPost handled inbound post requests
func httpPost(w http.ResponseWriter, r *http.Request) {
	endpoint := &Endpoint{}
	// Parse content body as endpoint struct
	err := json.NewDecoder(r.Body).Decode(endpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Insert the struct into the database
	err = insert(endpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	// Confirm the action was successful
	w.WriteHeader(http.StatusOK)
}

// httpPost handled inbound delete requests
func httpDelete(w http.ResponseWriter, r *http.Request) {
	endpoint := &Endpoint{}
	// Parse content body as endpoint struct
	err := json.NewDecoder(r.Body).Decode(endpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Insert the struct into the database
	err = drop(endpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	// Confirm the action was successful
	w.WriteHeader(http.StatusOK)
}

func httpPut(w http.ResponseWriter, r *http.Request) {
	endpoint := &Endpoint{}
	// Parse content body as endpoint struct
	err := json.NewDecoder(r.Body).Decode(endpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Insert the struct into the database
	err = update(endpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	// Confirm the action was successful
	w.WriteHeader(http.StatusOK)
}
