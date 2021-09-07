package server

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
)

func New(models ...interface{}) (router chi.Router, err error) {
	// Generate a new Mux
	r := chi.NewRouter()

	// Establish a database connection
	database, err := NewDatabase()
	if err != nil {
		return r, err
	}

	// Default Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Custom Middleware
	r.Use(databaseContext(database))
	r.Use(corsHeaders())

	// Emplace and generate ambiguous tables
	for _, model := range models {
		database.AutoMigrate(model)
	}

	return r, err
}

func ProtectedRoutes(r chi.Router, f func(r chi.Router)) {
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))

		r.Use(jwtauth.Authenticator)

		r.Group(f)
	})
}

func corsHeaders() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
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
