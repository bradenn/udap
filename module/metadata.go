// Copyright (c) 2021 Braden Nicholson

package module

type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
}

type UpdateBuffer struct {
	InstanceId string
	Data       string
}
