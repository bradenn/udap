package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type module struct{}

var Module module

type Api struct {
	Host string `json:"url"`
	Path string `json:"path"`
}

func (m *module) Run(context []byte, payload interface{}) (result interface{}, err error) {
	api := Api{}

	err = json.Unmarshal(context, &api)
	if err != nil {
		return nil, err
	}

	buffer, err := api.Get(payload)
	if err != nil {
		return nil, err
	}

	return buffer.String(), err

}

func (a *Api) Get(model interface{}) (buffer bytes.Buffer, err error) {
	url := fmt.Sprintf("%s%s", a.Host, a.Path)
	response, err := http.Get(url)
	if err != nil {
		return buffer, err
	}

	_, err = buffer.ReadFrom(response.Body)
	if err != nil {
		return buffer, err
	}

	return buffer, nil
}
