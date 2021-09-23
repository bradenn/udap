package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"net/http"
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
	err = connect()
	if err != nil {
		return Server{}, err
	}
	// Default Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	// Status Middleware
	router.Use(middleware.Heartbeat("/status"))
	// Custom Middleware
	router.Use(corsHeaders())
	// Seek, verify and validate JWT tokens
	router.Use(VerifyToken())
	return Server{router: router}, nil
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
	err := http.ListenAndServe("0.0.0.0:3020", s.router)
	if err != nil {
		return err
	}
	return nil
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
