package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"log"
	"net/http"
	"udap/config"
	"udap/module"
	"udap/server"
	"udap/server/auth"
)

func main() {
	var err error
	config.Init()

	router, err := server.New(Endpoint{}, Function{}, Entity{}, Group{})
	if err != nil {
		log.Fatalln(err)
	}

	// Routes required a valid JWT token
	router.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens

		r.Use(auth.VerifyToken())
		r.Use(auth.RequireAuth)

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["id"])))
		})
	})

	// Unsecured, public routes
	router.Group(func(r chi.Router) {

	})

	router.Route("/modules", func(r chi.Router) {
		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			req, _ := server.NewRequest(writer, request)

			type A struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}

			var model []A

			list, err := module.List()
			if err != nil {
				return
			}

			for _, m := range list {
				model = append(model, A{
					Name:        m.Name(),
					Description: m.Description(),
				})
			}

			req.JSON(model, http.StatusOK)
		})
	})

	router.Route("/groups", func(r chi.Router) {
		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			req, db := server.NewRequest(writer, request)

			var model []Group

			db.Model(&model).Find(&model)

			req.JSON(model, http.StatusOK)
		})
	})

	router.Route("/endpoints", func(r chi.Router) {

		end := Endpoint{}
		r.Post("/", end.Create)

		r.Group(func(t chi.Router) {

			r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
				req, db := server.NewRequest(writer, request)

				var model []Endpoint

				db.Model(&model).Find(&model)

				req.JSON(model, http.StatusOK)
			})
			r.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
				req, db := server.NewRequest(writer, request)

				var model Endpoint

				id := req.Param("id")

				db.Model(&model).Preload("Groups").Where("id = ?", id).First(&model)

				req.JSON(model, http.StatusOK)
			})
		})
	})

	err = http.ListenAndServe("0.0.0.0:3020", router)
	if err != nil {
		log.Fatalln(err)
	}

}
