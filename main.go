package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"udap/config"
	"udap/module"
	"udap/server"
)

func main() {
	var err error
	config.Init()

	a := server.WriteJwt(map[string]interface{}{"hello": "hi"})
	fmt.Println(a)

	p := server.CheckJwt(a)

	fmt.Println(p)

	s, err := server.New(Endpoint{}, Function{}, Entity{}, Group{})
	if err != nil {
		log.Fatalln(err)
	}

	s.Route("/auth", func(r chi.Router) {
		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			req, _ := server.NewRequest(writer, request)

			req.JSON("", http.StatusOK)
		})
	})

	server.ProtectedRoutes(s, func(r chi.Router) {
		r.Get("/status", func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("pong"))
		})
	})

	s.Route("/modules", func(r chi.Router) {
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

	s.Route("/groups", func(r chi.Router) {
		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			req, db := server.NewRequest(writer, request)

			var model []Group

			db.Model(&model).Find(&model)

			req.JSON(model, http.StatusOK)
		})
	})

	s.Route("/endpoints", func(r chi.Router) {
		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			req, db := server.NewRequest(writer, request)

			var model []Endpoint

			db.Model(&model).Preload("Groups").Find(&model)

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

	err = http.ListenAndServe("0.0.0.0:3020", s)
	if err != nil {
		log.Fatalln(err)
	}

}
