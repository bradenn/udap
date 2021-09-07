package main

import (
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

	s, err := server.New()
	if err != nil {
		log.Fatalln(err)
	}

	server.RegisterModel(Group{})
	server.RegisterModel(Function{})
	server.RegisterModel(Entity{})
	server.RegisterModel(Endpoint{})

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

	// 	// r.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
	// 	// 	db := server.dbContext(request)
	// 	// 	var function types.Function
	// 	//
	// 	// 	err = db.Model(&types.Function{}).Find(&function, "id = ?",
	// 	// 		chi.URLParam(request,
	// 	// 			"id")).Error
	// 	//
	// 	// 	if err != nil {
	// 	// 		fmt.Println(err)
	// 	// 	}
	// 	//
	// 	// 	server.JSON(function.Run(""), writer)
	// 	// })
	// })
	// s.Route("/endpoints", func(r chi.Router) {
	// 	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
	// 		db := server.dbContext(request)
	// 		var endpoints []Endpoint
	//
	// 		err = db.Model(&Endpoint{}).Find(&endpoints).Error
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	//
	// 		server.JSON(endpoints, writer)
	// 	})
	// 	r.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
	// 		db := server.dbContext(request)
	// 		var endpoint Endpoint
	//
	// 		err = db.Model(&Endpoint{}).Preload("Groups.Entities.Functions").Find(&endpoint, "id = ?",
	// 			chi.URLParam(request,
	// 				"id")).Error
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	//
	// 		server.JSON(endpoint, writer)
	// 	})
	// 	r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
	// 		db := server.dbContext(request)
	// 		endpoint := Endpoint{}
	//
	// 		err = server.DecodeJson(request, &endpoint)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	//
	// 		err = db.Create(&endpoint).Error
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	//
	// 		server.JSON(endpoint, writer)
	// 	})
	// })
	// s.Route("/entities", func(r chi.Router) {
	// 	r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
	// 		db := server.dbContext(request)
	// 		entity := Entity{}
	//
	// 		err = server.DecodeJson(request, &entity)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	//
	// 		err = db.Create(&entity).Error
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	//
	// 		marshal, err := json.Marshal(entity)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	//
	// 		_, err = writer.Write(marshal)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	//
	// 	})
	//
	// 	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
	// 		db := server.dbContext(request)
	// 		var entities []Entity
	//
	// 		err = db.Model(&Entity{}).First(&entities).Error
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	//
	// 		server.JSON(entities, writer)
	// 	})
	//
	// 	r.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
	// 		db := server.dbContext(request)
	// 		var entity Entity
	// 		id := chi.URLParam(request, "id")
	//
	// 		db.Model(&Entity{})
	// 		db.Preload("Functions").Preload("Groups")
	//
	// 		err = db.Where("id = ?", id).First(&entity).Error
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	//
	// 		server.JSON(entity, writer)
	// 	})
	// })

	err = http.ListenAndServe("0.0.0.0:3020", s)
	if err != nil {
		log.Fatalln(err)
	}

}
