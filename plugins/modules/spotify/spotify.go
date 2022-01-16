// Copyright (c) 2022 Braden Nicholson

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"net/url"
	"os"
	"time"
	"udap/internal/log"
	"udap/internal/models"
	"udap/pkg/plugin"
)

var Module Spotify

type Spotify struct {
	plugin.Module
	api      SpotifyApi
	Accounts map[string]SpotifyApi
	id       string
	local    MediaState
}

type Track struct {
}

type MediaState struct {
	Playing bool `json:"playing"`
	Data    json.RawMessage
	Current Track
}

func init() {
	config := plugin.Config{
		Name:        "spotify",
		Type:        "module",
		Description: "Single instance spotify controller",
		Version:     "2.0.1",
		Author:      "Braden Nicholson",
	}
	Module.Config = config
}

func (s *Spotify) PutAttribute(key string) models.FuncPut {
	return func(a models.Attribute) error {
		switch key {
		case "current":
			break
		}
		return nil
	}
}

func (s *Spotify) GetAttribute(key string) models.FuncGet {
	return func(a models.Attribute) (string, error) {
		switch key {
		case "current":
			return s.api.CurrentSong()
		}
		return "", nil
	}
}

func (s *Spotify) Setup() (plugin.Config, error) {
	e := models.NewMediaEntity("Remote", "spotify")

	_, err := s.Entities.Register(e)
	if err != nil {
		return plugin.Config{}, err
	}

	playing := models.Attribute{
		Key:    "current",
		Value:  "false",
		Type:   "toggle",
		Entity: e.Id,
	}
	playing.FnGet(s.GetAttribute(playing.Key))
	playing.FnPut(s.PutAttribute(playing.Key))

	err = s.Attributes.Register(&playing)
	if err != nil {
		return plugin.Config{}, err
	}

	if e.Id != "" {
		s.id = e.Id
		if e.Config == "" {
			a := SpotifyApi{}
			m, _ := json.Marshal(&a)
			e.Config = string(m)
			return s.Config, nil
		}
		a := SpotifyApi{}
		err = json.Unmarshal([]byte(e.Config), &a)
		if err != nil {
			return plugin.Config{}, err
		}
		s.api = a

	}
	return s.Config, nil
}

func (s *Spotify) Update() error {
	return nil
}

func (s *Spotify) Run() error {
	s.api.Authenticate()
	marshal, err := json.Marshal(s.api)
	if err != nil {
		return err
	}
	_, err = s.Entities.Config(s.id, string(marshal))
	if err != nil {
		return err
	}
	return nil
}

func (s Spotify) Poll(v string) (string, error) {
	spotifyApi := SpotifyApi{}
	err := json.Unmarshal([]byte(v), &spotifyApi)
	if err != nil {
		return "", err
	}
	return spotifyApi.CurrentSong()
}

func (s *SpotifyApi) CurrentSong() (string, error) {
	return s.authenticatedRequest("GET", "me/player/currently-playing")
}

func (s *SpotifyApi) Play(_ string) (string, error) {
	return s.authenticatedRequest("PUT", "me/player/play")
}

func (s *SpotifyApi) Pause(_ string) (string, error) {
	return s.authenticatedRequest("PUT", "me/player/pause")
}

func (s *SpotifyApi) authenticatedRequest(method string, path string) (string, error) {
	access, _ := s.GetToken()
	urlString := fmt.Sprintf("https://api.spotify.com/v1/%s", path)
	request, _ := http.NewRequest(method, urlString, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type SpotifyCallback struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

const (
	CALLBACK = "http://localhost:9966/callback"
)

type SpotifyApi struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	Refresh   string    `json:"refresh"`
	Data      json.RawMessage
}

func (s *SpotifyApi) Authenticate() {
	if s.Token == "" {
		log.Critical("ACTION REQUIRED: %s", s.loginURL())
		done := make(chan bool)
		e := make(chan string)
		go func() {
			defer func() {
				done <- true
			}()
			s.beginListening(e)
		}()
		msg := <-e
		s.requestToken(msg)
		<-done
	}
}

func (s *SpotifyApi) GetToken() (string, error) {
	if s.Token == "" {
		return "", fmt.Errorf("not authenticated")
	} else if time.Now().After(s.ExpiresAt) {
		s.refreshAccess()
	}
	return s.Token, nil
}

func (s *SpotifyApi) beginListening(e chan string) {

	r := chi.NewRouter()
	srv := http.Server{
		Addr: "0.0.0.0:9966",
	}

	r.Get("/callback", func(w http.ResponseWriter, req *http.Request) {
		code := req.URL.Query().Get("code")
		e <- code

		log.Event("Spotify: Authentication succeeded")

		_, err := w.Write([]byte("All done!"))
		if err != nil {
			return
		}
		err = srv.Close()
		if err != nil {
			return
		}
	})

	srv.Handler = r
	err := srv.ListenAndServe()
	if err != nil {
		return
	}
}

func (s *SpotifyApi) loginURL() string {

	scope := "user-modify-playback-state user-read-currently-playing"
	clientId := os.Getenv("spotifyClient")

	urlString := fmt.Sprintf("https://accounts.spotify."+
		"com/authorize?response_type=%s&client_id=%s&scope=%s&redirect_uri=%s",
		"code", clientId, url.PathEscape(scope), url.PathEscape(CALLBACK))

	return urlString
}

func (s *SpotifyApi) requestToken(code string) {
	values := url.Values{}

	values.Set("grant_type", "authorization_code")
	values.Set("code", code)
	values.Set("redirect_uri", CALLBACK)

	request, err := s.basicRequest("https://accounts.spotify.com/api/token", values)
	if err != nil {
		return
	}

	callbackBody := SpotifyCallback{}
	err = json.Unmarshal([]byte(request), &callbackBody)
	if err != nil {
		return
	}
	s.Token = callbackBody.AccessToken
	s.Refresh = callbackBody.RefreshToken
	s.ExpiresAt = time.Now().Add(time.Second * time.Duration(callbackBody.ExpiresIn))
}

func (s *SpotifyApi) refreshAccess() {
	values := url.Values{}

	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", s.Refresh)
	request, err := s.basicRequest("https://accounts.spotify.com/api/token", values)
	if err != nil {
		return
	}

	refresh := SpotifyRefresh{}
	err = json.Unmarshal([]byte(request), &refresh)
	if err != nil {
		return
	}
	s.Token = refresh.AccessToken
	s.ExpiresAt = time.Now().Add(time.Second * time.Duration(refresh.ExpiresIn))

}

type SpotifyRefresh struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
}

func (s *SpotifyApi) basicRequest(path string, values url.Values) (string, error) {

	id := os.Getenv("spotifyClient")
	secret := os.Getenv("spotifySecret")

	base := fmt.Sprintf("%s:%s", id, secret)
	encoded := base64.StdEncoding.EncodeToString([]byte(base))

	var buf bytes.Buffer
	body := values.Encode()
	buf.Write([]byte(body))

	request, err := http.NewRequest("POST", path, &buf)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", encoded))

	request.PostForm = values
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	var res bytes.Buffer

	_, err = res.ReadFrom(response.Body)
	if err != nil {
		return res.String(), err
	}

	return res.String(), nil
}
