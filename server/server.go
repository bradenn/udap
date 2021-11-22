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
	"udap/udap"
	"udap/udap/db"
)

type Routable interface {
	Route(router chi.Router)
}

type Server struct {
	router chi.Router
}

func New() (s Server, err error) {
	// Load Authenticate
	Init()
	// Generate a new Mux
	router := chi.NewRouter()
	// Establish a database connection
	// Default Middleware
	router.Use(udap.Middleware)
	router.Use(middleware.Recoverer)
	// Custom Middleware
	router.Use(corsHeaders())
	// Status Middleware
	router.Use(middleware.Heartbeat("/status"))
	// Seek, verify and validate JWT tokens
	router.Use(VerifyToken())

	return Server{router: router}, nil
}

func databaseContext(database *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "DB", database)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (s *Server) Migrate(inf interface{}) {
	err := db.DB.AutoMigrate(inf)
	if err != nil {
		udap.Error(err.Error())
	}
}

func (s *Server) RouteSecure(path string, handler func(r chi.Router)) {
	s.router.Group(func(r chi.Router) {
		// Enforce tokens
		r.Use(RequireAuth)
		// Begin integration of authorized routes
		r.Route(path, handler)
	})
}

func (s *Server) RoutePublic(path string, handler func(r chi.Router)) {
	s.router.Group(func(r chi.Router) {
		// Begin integration of unauthorized routes
		r.Route(path, handler)
	})
}

func (s *Server) Router() chi.Router {
	return s.router
}

func (s *Server) Run() error {
	host := fmt.Sprintf("%s:%s", os.Getenv("hostname"), os.Getenv("port"))
	err := http.ListenAndServe(host, s.router)
	if err != nil {
		return err
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
