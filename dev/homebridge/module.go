package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"udap/template"
)

var Export HomeBridge

type Environment struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (e Environment) toString() string {
	marshal, _ := json.Marshal(e)
	return string(marshal)
}

type Instance struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (i Instance) toString() string {
	marshal, _ := json.Marshal(i)
	return string(marshal)
}

func (h *HomeBridge) Authenticate() (string, error) {
	var buf bytes.Buffer

	s := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		h.Environment.Username,
		h.Environment.Password,
	}

	marshal, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	_, err = buf.Write(marshal)
	if err != nil {
		return "", err
	}

	host := h.Environment.Host
	port := h.Environment.Port

	path := fmt.Sprintf("http://%s:%s/api/auth/login", host, port)

	request, err := http.NewRequest("POST", path, &buf)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	var res bytes.Buffer

	_, err = res.ReadFrom(response.Body)
	if err != nil {
		return "", err

	}

	i := Instance{}

	err = json.Unmarshal(res.Bytes(), &i)
	if err != nil {
		return "", err
	}

	return res.String(), nil
}

type HomeBridge struct {
	Environment Environment
}

func (h *HomeBridge) InitInstance() (string, error) {
	return h.Authenticate()
}

func (h *HomeBridge) Initialize(env string) {
	environment := Environment{}
	err := json.Unmarshal([]byte(env), &environment)
	if err != nil {
		return
	}
	h.Environment = environment
}

func (h HomeBridge) Metadata() template.Metadata {
	metadata := template.Metadata{
		Name:        "Homebridge",
		Description: "Connect to Homebridge, and transitively, homekit.",
		Version:     "1.0.0",
		Author:      "Braden Nicholson",
	}
	return metadata
}

func (h HomeBridge) Poll(v string) (string, error) {
	// path := fmt.Sprintf("ws://%s:%s/socket.io/?token=%s&EIO=3&transport=websocket")
	// websocket.Dial("ws:/")
	return "", nil
}

func (h HomeBridge) Run(v string, action string) (string, error) {
	return "", nil
}
