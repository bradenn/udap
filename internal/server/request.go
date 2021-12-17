// Copyright (c) 2021 Braden Nicholson

package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
	"net/http"
)

type Req struct {
	writer  http.ResponseWriter
	request *http.Request
}

type RejectError struct {
	Error interface{} `json:"error"`
}

func NewRequest(writer http.ResponseWriter, request *http.Request) (req *Req, db *gorm.DB) {
	req = &Req{
		writer:  writer,
		request: request,
	}

	db = request.Context().Value("DB").(*gorm.DB)
	return req, db
}

func (r *Req) JWTClaim(key string) interface{} {
	_, claims, err := jwtauth.FromContext(r.request.Context())

	if err != nil {
		return ""
	}

	return claims[key]
}

func (r *Req) Param(key string) string {
	return chi.URLParam(r.request, key)
}

func (r *Req) ParamObjectId(key string) (primitive.ObjectID, error) {
	id := r.Param(key)
	if !primitive.IsValidObjectID(id) {
		return primitive.ObjectID{}, fmt.Errorf("invalid objectId")
	}
	return primitive.ObjectIDFromHex(id)
}

func (r *Req) Body() string {
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(r.request.Body)
	if err != nil {
		return ""
	}

	return buffer.String()
}

func (r *Req) DecodeModel(model interface{}) {
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(r.request.Body)
	if err != nil {
		r.Reject(err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(buffer.Bytes(), model)
	if err != nil {
		r.Reject(err.Error(), http.StatusBadRequest)
	}
}

func (r *Req) Reject(payload interface{}, status int) {

	errPayload := RejectError{Error: payload}

	marshal, err := json.Marshal(errPayload)
	if err != nil {
		fmt.Println(err)
	}

	writeCors(r.writer)

	r.writer.Header().Set("Content-Type", "application/json")
	r.writer.WriteHeader(status)

	_, err = r.writer.Write(marshal)
	if err != nil {
		return
	}
}

func (r *Req) ResolveRaw(payload string, status int) {
	writeCors(r.writer)

	r.writer.Header().Set("Content-Type", "application/json")
	r.writer.WriteHeader(status)

	_, err := r.writer.Write([]byte(payload))
	if err != nil {
		return
	}
}

func (r *Req) Resolve(payload interface{}, status int) {
	marshal, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}

	writeCors(r.writer)

	r.writer.Header().Set("Content-Type", "application/json")
	r.writer.WriteHeader(status)

	_, err = r.writer.Write(marshal)
	if err != nil {
		return
	}
}

func writeCors(writer http.ResponseWriter) {
	writer.Header().Set("Access-Control-Allow-Headers", "X-API-KEY, Origin, X-Requested-With, Content-Type, Accept, Access-Control-Request-Method")
}
