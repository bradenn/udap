package main

import (
	"time"
	"udap/template"
)

var Export template.Module

func init() {

	functions := map[string]template.Function{}

	functions["localTime"] = LocalTime

	metadata := template.Metadata{
		Name:        "Timekeeper",
		Description: "Get various time related information.",
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

func LocalTime(_ string) (string, error) {
	return time.Now().String(), nil
}
