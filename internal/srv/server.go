// Copyright (c) 2022 Braden Nicholson

package srv

import (
	"github.com/go-chi/chi"
	"net/http"
	"time"
	"udap/internal/port/routes"
	"udap/platform/jwt"
	"udap/platform/router"
)

type Server struct {
	router chi.Router
	server *http.Server
}

func NewServer() Server {
	srv := Server{}
	srv.router = router.New()

	srv.server = &http.Server{
		Addr:              ":3020",
		Handler:           srv.router,
		ReadTimeout:       time.Second,
		WriteTimeout:      time.Second * 2,
		IdleTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 2,
	}

	return srv
}

func (s *Server) AddRoute(route routes.Routable) {
	s.router.Group(func(internal chi.Router) {
		internal.Use(jwt.Authenticator)
		route.RouteInternal(internal)
	})
	route.RouteExternal(s.router)
}

func (s *Server) AddRoutes(routable ...routes.Routable) {
	for _, route := range routable {
		s.router.Group(func(internal chi.Router) {
			internal.Use(jwt.Authenticator)
			route.RouteInternal(internal)
		})
		route.RouteExternal(s.router)
	}
}

func (s *Server) Run() error {
	err := s.server.ListenAndServeTLS("./certs/ca.crt", "./certs/ca.key")
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Close() error {
	err := s.server.Close()
	if err != nil {
		return err
	}
	return nil
}
