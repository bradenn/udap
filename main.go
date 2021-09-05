package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
	"udap/config"
	"udap/server"
)

type Persistent struct {
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" sql:"index"`
	Id        uuid.UUID  `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
}

type Endpoint struct {
	Persistent
	Name   string  `json:"name" gorm:"unique"`
	Groups []Group `json:"groups" gorm:"many2many:endpointGroup;"`
}

type Entity struct {
	Persistent
	Name        string     `json:"name" gorm:"unique"`
	Description string     `json:"description"`
	Functions   []Function `json:"functions" gorm:"many2many:entityFunction;"`
}

type Group struct {
	Persistent
	Name       string   `json:"name"  gorm:"unique"`
	Entities   []Entity `json:"entities" gorm:"many2many:entityGroup;"`
	Identifier string   `json:"identifier"  gorm:"unique"`
}

func JSON(payload interface{}, writer http.ResponseWriter) {

	marshal, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}

	writer.Header().Set("Access-Control-Allow-Headers", "X-API-KEY, Origin, X-Requested-With, Content-Type, Accept, Access-Control-Request-Method")
	writer.Header().Set("Content-Type", "json")
	writer.WriteHeader(200)

	_, err = writer.Write(marshal)
	if err != nil {
		return
	}
}

func main() {
	var err error
	config.Init()

	s, err := server.New()
	if err != nil {
		log.Fatalln(err)
	}

	server.RegisterType(Group{})
	server.RegisterType(Function{})
	server.RegisterType(Entity{})
	server.RegisterType(Endpoint{})

	s.Route("/function", func(r chi.Router) {
		r.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
			db := server.DbContext(request)
			var function Function

			err = db.Model(&Function{}).Find(&function, "id = ?",
				chi.URLParam(request,
					"id")).Error

			if err != nil {
				fmt.Println(err)
			}

			JSON(function.Run(""), writer)
		})
	})
	s.Route("/endpoints", func(r chi.Router) {
		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			db := server.DbContext(request)
			var endpoints []Endpoint

			err = db.Model(&Endpoint{}).Find(&endpoints).Error
			if err != nil {
				fmt.Println(err)
			}

			JSON(endpoints, writer)
		})
		r.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
			db := server.DbContext(request)
			var endpoint Endpoint

			err = db.Model(&Endpoint{}).Preload("Groups.Entities.Functions").Find(&endpoint, "id = ?",
				chi.URLParam(request,
					"id")).Error
			if err != nil {
				fmt.Println(err)
			}

			JSON(endpoint, writer)
		})
		r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
			db := server.DbContext(request)
			endpoint := Endpoint{}

			err = server.DecodeJson(request, &endpoint)
			if err != nil {
				fmt.Println(err)
			}

			err = db.Create(&endpoint).Error
			if err != nil {
				fmt.Println(err)
			}

			JSON(endpoint, writer)
		})
	})
	s.Route("/entities", func(r chi.Router) {
		r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
			db := server.DbContext(request)
			entity := Entity{}

			err = server.DecodeJson(request, &entity)
			if err != nil {
				fmt.Println(err)
			}

			err = db.Create(&entity).Error
			if err != nil {
				fmt.Println(err)
			}

			marshal, err := json.Marshal(entity)
			if err != nil {
				fmt.Println(err)
			}

			_, err = writer.Write(marshal)
			if err != nil {
				fmt.Println(err)
			}

		})

		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			db := server.DbContext(request)
			var entities []Entity

			err = db.Model(&Entity{}).First(&entities).Error
			if err != nil {
				fmt.Println(err)
			}

			marshal, err := json.Marshal(entities)
			if err != nil {
				fmt.Println(err)
			}
			writer.Header().Set("Access-Control-Allow-Headers", "X-API-KEY, Origin, X-Requested-With, Content-Type, Accept, Access-Control-Request-Method")
			writer.Header().Set("Content-Type", "json")
			writer.WriteHeader(200)
			writer.Write(marshal)

		})

		r.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
			db := server.DbContext(request)
			var entity Entity

			err = db.Model(&Entity{}).Preload("Functions").Preload("Groups").Where("id = ?",
				chi.URLParam(request, "id")).First(&entity).Error
			if err != nil {
				fmt.Println(err)
			}

			marshal, err := json.Marshal(entity)
			if err != nil {
				fmt.Println(err)
			}
			writer.Header().Set("Access-Control-Allow-Headers", "X-API-KEY, Origin, X-Requested-With, Content-Type, Accept, Access-Control-Request-Method")
			writer.Header().Set("Content-Type", "json")
			writer.WriteHeader(200)
			writer.Write(marshal)

		})
	})

	err = http.ListenAndServe("0.0.0.0:3020", s)
	if err != nil {
		log.Fatalln(err)
	}

}
