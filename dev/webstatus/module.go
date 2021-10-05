package main

import (
	"encoding/json"
	"net/http"
	"udap/template"
)

// Export Allows the parent process to read
var Export WebStatus

func init() {
	Export = WebStatus{}
}

type WebStatus struct {
}

type Configuration struct {
	Url string `json:"url"`
}

func (w *WebStatus) InitInstance() (string, error) {
	sample := &Configuration{
		Url: "https://google.com",
	}
	env, err := json.Marshal(sample)
	if err != nil {
		return "", err
	}
	return string(env), nil
}

func (w *WebStatus) Initialize(env string) {

}

func (w *WebStatus) Metadata() template.Metadata {
	return template.Metadata{
		Name:        "Web Status",
		Description: "Determine whether a website is down or not",
		Version:     "1.0.0",
		Author:      "Braden Nicholson",
	}
}

func (w *WebStatus) Poll(v string) (string, error) {
	conf := Configuration{}
	err := json.Unmarshal([]byte(v), &conf)
	if err != nil {
		return "", err
	}

	get, err := http.Get(conf.Url)
	if err != nil {
		return err.Error(), err
	}
	returnStruct := struct {
		Url    string `json:"url"`
		Status string `json:"status"`
	}{
		Url:    conf.Url,
		Status: get.Status,
	}

	marshal, err := json.Marshal(returnStruct)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

func (w *WebStatus) Run(v string, action string) (string, error) {
	return "", nil
}
