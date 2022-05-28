// Copyright (c) 2022 Braden Nicholson

package orchestrator

import (
	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"net/http"
	"udap/internal/controller"
	"udap/internal/log"
	"udap/internal/port/routes"
	"udap/internal/port/runtimes"
	"udap/platform/database"
	"udap/platform/router"
)

type orchestrator struct {
	db         *gorm.DB
	router     chi.Router
	controller *controller.Controller
}

func (o orchestrator) Run() error {

	server := &http.Server{Addr: ":8080", Handler: o.router}

	err := server.ListenAndServe()
	if err != nil {
		log.ErrF(err, "http server exited with error:\n")
	}

	return nil
}

type Orchestrator interface {
	Init() error
	Run() error
}

func NewOrchestrator() Orchestrator {
	// Initialize Database
	db, err := database.New()
	if err != nil {
		return nil
	}
	// Initialize Router
	r := router.New()

	return &orchestrator{
		db:         db,
		router:     r,
		controller: nil,
	}
}

func (o orchestrator) Init() error {
	var err error
	o.controller, err = controller.NewController(o.db)
	if err != nil {
		return err
	}

	// Initialize and route applicable domains
	routes.NewUserRouter(o.controller.Users).RouteUsers(o.router)
	routes.NewEndpointRouter(o.controller.Endpoints).RouteEndpoints(o.router)
	routes.NewModuleRouter(o.controller.Modules).RouteModules(o.router)

	runtimes.NewModuleRuntime(o.controller.Modules)

	return nil
}
