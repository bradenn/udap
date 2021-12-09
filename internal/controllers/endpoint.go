// Copyright (c) 2021 Braden Nicholson

package controllers

import (
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/models"
	"udap/internal/server"
	"udap/internal/store"
)

//
// type endpointSocket struct{}
//
// func (e *endpointSocket) create(request server.Request) {
// 	// body := request.Body()
// 	// Make sure the name is unique
// 	// Make sure the endpoint callers is a allowed
// 	// Make the endpoint
// 	// Return the endpoint join code
//
// 	request.Resolve(models.Endpoint{})
// }
//
// func (e *endpointSocket) delete(request server.Request) {
//
// }
//
// func (e *endpointSocket) update(request server.Request) {
//
// }
//
// func EndpointSocket(s server.Socket) {
// 	ec := endpointSocket{}
// 	s.Handle("endpoint/create", ec.create)
// 	s.Handle("endpoint/{id}/delete", ec.delete)
// 	s.Handle("endpoint/{id}/update", ec.update)
// }

type EndpointController struct{}

func EndpointRouter(router chi.Router) {
	ec := EndpointController{}
	router.Get("/register/{accessKey}", ec.register)
}

func (e *EndpointController) register(w http.ResponseWriter, r *http.Request) {
	req, _ := server.NewRequest(w, r)
	key := chi.URLParam(r, "accessKey")
	endpoint := models.Endpoint{}

	err := store.DB.Model(&models.Endpoint{}).Where("key = ?", key).First(&endpoint).Error
	if err != nil {
		req.Reject(err, http.StatusBadRequest)
	}

	jwt, err := server.SignUUID(endpoint.Id)
	if err != nil {
		req.Reject("Internal Error", http.StatusInternalServerError)
		return
	}

	resolve := map[string]interface{}{"token": jwt}

	req.Resolve(resolve, http.StatusOK)
}
