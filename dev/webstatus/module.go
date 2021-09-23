package main

import (
	"net/http"
	"udap/template"
)

var Export template.Module

func init() {

	functions := map[string]template.Function{}

	functions["isDown"] = IsDown

	metadata := template.Metadata{
		Name:        "Web Status",
		Description: "Determine whether a website is down or not",
		Version:     "1.0.0",
		Author:      "Braden Nicholson",
	}

	module := template.NewModule(metadata, functions, Configure)

	Export = module

}

func Configure() {

	// config := Export.GetConfig()
	//
	// instance := Export.GetInstance().String()

}

func IsDown(url string) (string, error) {
	get, err := http.Get(url)
	if err != nil {
		return "Down", err
	}

	if get.StatusCode != http.StatusOK {
		return get.Status, err
	}

	return get.Status, nil
}
