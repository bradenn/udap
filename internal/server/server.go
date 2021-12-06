// Copyright (c) 2021 Braden Nicholson

package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"os"
)

type Server struct {
	chi.Router
	host string
}

// Dependency is the level at which this service needs to run
func (w *Server) Dependency() (level int) {
	return 2
}
func (w *Server) Name() (name string) {
	return "server"
}

// Load is a blocking server function
func (w *Server) Load() (err error) {
	w.Router = chi.NewRouter()

	w.Use(middleware.Recoverer)
	// Custom Middleware
	w.Use(corsHeaders())
	// Status Middleware
	w.Use(middleware.Heartbeat("/status"))
	// Seek, verify and validate JWT tokens
	w.Use(verifyToken())
	// Load JWT Keys
	LoadKeys()

	// Format connection string
	w.host = fmt.Sprintf("%s:%s", os.Getenv("hostname"), os.Getenv("port"))
	return nil
}

// Run is a blocking server function
func (w *Server) Run(ctx context.Context) (err error) {
	// Listen and server indefinitely
	err = http.ListenAndServe(w.host, w)
	if err != nil {
		return
	}
	return nil
}

func corsHeaders() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}

func (w *Server) RouteSecure(path string, handler func(r chi.Router)) {
	w.Group(func(r chi.Router) {
		// Enforce tokens
		r.Use(RequireAuth)
		// Begin integration of authorized routes
		r.Route(path, handler)
	})
}

func (w *Server) RoutePublic(path string, handler func(r chi.Router)) {
	w.Group(func(r chi.Router) {
		// Begin integration of unauthorized routes
		r.Route(path, handler)
	})
}
