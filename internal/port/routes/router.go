// Copyright (c) 2022 Braden Nicholson

package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"io"
	"net/http"
)

type Routable interface {
	RouteInternal(router chi.Router)
	RouteExternal(router chi.Router)
}

type PersistentRoute interface {
	create(w http.ResponseWriter, req *http.Request)
	update(w http.ResponseWriter, req *http.Request)
	delete(w http.ResponseWriter, req *http.Request)
}

type ServiceInterface[T any] interface {
	FindAll() (*[]T, error)
	FindById(id string) (*T, error)
	Create(t *T) error
	Update(t *T) error
	Delete(id string) error
}

type pRoute[T any] struct {
	ServiceInterface[T]
	Container T
}

func (p pRoute[T]) create(w http.ResponseWriter, req *http.Request) {
	buf := bytes.Buffer{}

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(w, "could not read object", 500)
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(req.Body)

	sr := p.Container

	err = json.Unmarshal(buf.Bytes(), &sr)
	if err != nil {
		http.Error(w, "could not parse object", 500)
		return
	}

	err = p.Create(&sr)
	if err != nil {
		http.Error(w, "could not create object", 500)
		return
	}

	_, err = w.Write([]byte("OK"))
	if err != nil {
		return
	}

}

func (p pRoute[T]) update(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "id")
	if key == "" {
		http.Error(w, "access key not provided", 401)
		return
	}

	err := p.Delete(key)
	if err != nil {
		http.Error(w, "could not update object", 500)
		return
	}

	_, err = w.Write([]byte("OK"))
	if err != nil {
		return
	}
}

func (p pRoute[T]) delete(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "id")
	if key == "" {
		http.Error(w, "access key not provided", 401)
		return
	}

	err := p.Delete(key)
	if err != nil {
		http.Error(w, "could not delete object", 500)
		return
	}

	_, err = w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
