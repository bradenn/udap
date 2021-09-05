package server

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
)

var database *gorm.DB

func New() (router chi.Router, err error) {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	database, err = NewDatabase()
	if err != nil {
		return r, err
	}

	// Inject Database context reference
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "DB", database)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	return r, err
}

func RegisterType(model interface{}) {
	database.AutoMigrate(model)
}

func DecodeJson(r *http.Request, model interface{}) error {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf.Bytes(), model)
	if err != nil {
		return err
	}

	return nil
}

func DbContext(r *http.Request) *gorm.DB {
	db := r.Context().Value("DB").(*gorm.DB)
	return db
}
