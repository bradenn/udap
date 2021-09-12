package server

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
)

type Routable interface {
	Route(router chi.Router)
}

type Server struct {
	router   chi.Router
	database *gorm.DB
}

func New() (s Server, err error) {
	// Load Authenticate
	Init()
	// Generate a new Mux
	router := chi.NewRouter()
	// Establish a database connection
	database, err := NewDatabase()
	if err != nil {
		return s, err
	}
	// Default Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	// Status Middleware
	router.Use(middleware.Heartbeat("/status"))
	// Custom Middleware
	router.Use(databaseContext(database))
	router.Use(corsHeaders())
	// Seek, verify and validate JWT tokens
	router.Use(VerifyToken())
	return Server{router: router, database: database}, nil
}

func (s *Server) RouteSecure(path string, routable Routable) {
	s.database.AutoMigrate(routable)
	s.router.Group(func(r chi.Router) {
		// Enforce tokens
		r.Use(RequireAuth)
		// Begin integration of authorized routes
		r.Route(path, routable.Route)
	})
}

func (s *Server) RoutePublic(path string, routable Routable) {
	s.database.AutoMigrate(routable)
	s.router.Group(func(r chi.Router) {
		// Begin integration of unauthorized routes
		r.Route(path, routable.Route)
	})
}

func (s *Server) Register(routable Routable) {
	s.database.AutoMigrate(routable)
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

func databaseContext(database *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "DB", database)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
