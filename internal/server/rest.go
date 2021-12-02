// Copyright (c) 2021 Braden Nicholson

package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"gorm.io/gorm"
	"net/http"
	"os"
	"udap/internal/log"
)

var server *Server

type Server struct {
	router chi.Router
	host   string
}

// Dependency is the level at which this service needs to run
func (w *Server) Dependency() (level int) {
	return 2
}
func (w *Server) Name() (name string) {
	return "webserver"
}

func (w *Server) Load() (err error) {
	server = w

	w.router = chi.NewRouter()

	w.router.Use(middleware.Recoverer)
	// Custom Middleware
	w.router.Use(corsHeaders())
	// Status Middleware
	w.router.Use(middleware.Heartbeat("/status"))
	// Seek, verify and validate JWT tokens
	w.router.Use(verifyToken())
	// Load JWT Keys
	LoadKeys()

	return nil
}

// Run is a blocking server function
func (w *Server) Run(interface{}) (err error) {
	// Format connection string
	w.host = fmt.Sprintf("%s:%s", os.Getenv("hostname"), os.Getenv("port"))
	log.Sherlock("Server Listening on '%s/socket/{token}'", w.host)
	// Listen and server indefinitely
	err = http.ListenAndServe(w.host, w.router)
	if err != nil {
		return err
	}
	return nil
}

func (w *Server) Cleanup() (err error) {
	return nil
}

func (w *Server) Router() chi.Router {
	return w.router
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

func databaseContext(database *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "DB", database)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (w *Server) RouteSecure(path string, handler func(r chi.Router)) {
	w.router.Group(func(r chi.Router) {
		// Enforce tokens
		r.Use(RequireAuth)
		// Begin integration of authorized routes
		r.Route(path, handler)
	})
}

func (w *Server) RoutePublic(path string, handler func(r chi.Router)) {
	w.router.Group(func(r chi.Router) {
		// Begin integration of unauthorized routes
		r.Route(path, handler)
	})
}

func (w *Server) socket() {

}
