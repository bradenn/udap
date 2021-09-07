package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Request struct {
	writer  http.ResponseWriter
	request *http.Request
}

func NewRequest(writer http.ResponseWriter, request *http.Request) (*Request, *gorm.DB) {
	req := &Request{
		writer:  writer,
		request: request,
	}
	db := dbContext(request)
	return req, db
}

func (r *Request) Param(key string) string {
	return chi.URLParam(r.request, key)
}

func (r *Request) Body() string {
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(r.request.Body)
	if err != nil {
		return ""
	}

	return buffer.String()

}

func (r *Request) DecodeModel(model interface{}) error {
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(r.request.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buffer.Bytes(), model)
	if err != nil {
		return err
	}

	return nil
}

func (r *Request) JSON(payload interface{}, status int) {
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

func dbContext(request *http.Request) *gorm.DB {
	db := request.Context().Value("DB").(*gorm.DB)
	return db
}
