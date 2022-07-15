// Copyright (c) 2022 Braden Nicholson

package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"udap/internal/plugin"
)

var Module Worldspace

type Worldspace struct {
	plugin.Module
	server http.Server
}

func init() {
	config := plugin.Config{
		Name:        "worldspace",
		Type:        "module",
		Description: "worldspace integration",
		Version:     "0.1.1",
		Author:      "Braden Nicholson",
	}
	Module.Config = config
}

func (w *Worldspace) Setup() (plugin.Config, error) {
	fmt.Println("Initializing worldspace")
	err := w.UpdateInterval(2000)
	if err != nil {
		return plugin.Config{}, err
	}
	return w.Config, nil
}

func (w *Worldspace) Update() error {
	if w.Ready() {

	}
	return nil
}

func (w *Worldspace) handleMotion(writer http.ResponseWriter, request *http.Request) {
	/*zone := chi.URLParam(request, "zone")*/

}

func (w *Worldspace) Run() error {
	fmt.Println("Starting worldspace")
	w.server = http.Server{}
	w.server.Addr = ":5055"
	router := chi.NewRouter()
	router.Post("/motion/{zone}", w.handleMotion)

	go func() {
		err := w.server.ListenAndServe()
		if err != nil {
			return
		}
	}()
	return nil
}

func (w *Worldspace) Dispose() error {
	fmt.Println("Disposing worldspace")
	ctx := context.Background()
	err := w.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
