// Copyright (c) 2022 Braden Nicholson

package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"udap/platform/jwt"
)

func New() chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	// Custom Middleware
	router.Use(corsHeaders())
	// Status Middleware
	router.Use(middleware.Heartbeat("/status"))
	// Seek, verify and validate JWT tokens
	router.Use(jwt.VerifyToken())
	// Load JWT Keys
	jwt.LoadKeys()
	return router
}

func corsHeaders() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Bond"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by interface{} of major browsers
	})
}
