package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
)

var database *gorm.DB

func New() (router chi.Router, err error) {
	r := chi.NewRouter()

	database, err = NewDatabase()
	if err != nil {
		return r, err
	}

	r.Use(databaseContext)
	r.Use(corsHeaders())

	return r, err
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

func routeModel(model interface{}, router chi.Router) {
	// registerModel(model)
	path := fmt.Sprintf("/%s", database.NewScope(model).TableName())
	router.Route(path, func(r chi.Router) {

		r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
			req, db := NewRequest(writer, request)
			err := req.DecodeModel(&model)
			if err != nil {
				req.JSON(model, http.StatusBadRequest)
			}

			err = db.Model(model).Create(model).Error
			if err != nil {
				req.JSON(err, http.StatusBadRequest)
				return
			}
			req.JSON(model, http.StatusOK)
		})
	})
}

func databaseContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "DB", database)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RegisterModel(model interface{}) {
	database.AutoMigrate(model)
}
