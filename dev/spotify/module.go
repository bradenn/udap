package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
	"udap/template"
)

var Export template.Module

var spotify SpotifyApi

func init() {

	functions := map[string]template.Function{}

	functions["currentlyPlaying"] = CurrentSong
	functions["pause"] = Pause
	functions["play"] = Play

	metadata := template.Metadata{
		Name:        "Spotify",
		Description: "Remote control for spotify",
		Version:     "1.0.1",
		Author:      "Braden Nicholson",
	}

	module := template.NewModule(metadata, functions, Configure)

	Export = module

}

func Configure() {

	config := Export.GetConfig()

	instance := Export.GetInstance().String()

	if !config.IsSet(instance) {
		spotify = SpotifyApi{}
		spotify.Authenticate()
		config.Set(instance, spotify)
	}

	err := config.Get(instance, &spotify)
	if err != nil {
		panic(err)
	}

}

func CurrentSong(_ string) (string, error) {
	request, err := authenticatedRequest("PUT", "me/player/currently-playing")
	if err != nil {
		return "", err
	}
	return request, nil
}

func Pause(_ string) (string, error) {
	request, err := authenticatedRequest("PUT", "me/player/pause")
	if err != nil {
		return "", err
	}
	return request, nil
}

func Play(_ string) (string, error) {
	request, err := authenticatedRequest("PUT", "me/player/play")
	if err != nil {
		return "", err
	}
	return request, nil
}

func authenticatedRequest(method string, path string) (string, error) {
	access, _ := spotify.GetToken()
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
}

func (a *SpotifyApi) Authenticate() {
	if a.Token == "" {
		fmt.Printf("Please log into your spotify account here: %s", a.loginURL())

		e := make(chan string)
		go a.beginListening(e)
		msg := <-e
		a.requestToken(msg)
	}
}

func (a *SpotifyApi) GetToken() (string, error) {
	if a.Token == "" {
		return "", fmt.Errorf("not authenticated")
	} else if time.Now().After(a.ExpiresAt) {
		a.refreshAccess()
	}
	return a.Token, nil
}

func (a *SpotifyApi) beginListening(e chan string) {

	srv := http.Server{
		Addr: "0.0.0.0:9966",
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		e <- code
		_, err := fmt.Fprintf(w, "Authenticated successfully... You can close this tab.")
		if err != nil {
			return
		}
		err = srv.Shutdown(context.Background())
		if err != nil {
			return
		}
	})

	err := srv.ListenAndServe()
	if err != nil {
		return
	}
}

func (a *SpotifyApi) loginURL() string {

	scope := "user-modify-playback-state user-read-currently-playing"
	clientId := os.Getenv("client")

	urlString := fmt.Sprintf("https://accounts.spotify."+
		"com/authorize?response_type=%s&client_id=%s&scope=%s&redirect_uri=%s",
		"code", clientId, url.PathEscape(scope), url.PathEscape(CALLBACK))

	return urlString
}

func (a *SpotifyApi) requestToken(code string) {
	values := url.Values{}

	values.Set("grant_type", "authorization_code")
	values.Set("code", code)
	values.Set("redirect_uri", CALLBACK)

	request, err := a.basicRequest("https://accounts.spotify.com/api/token", values)
	if err != nil {
		return
	}

	callbackBody := SpotifyCallback{}
	err = json.Unmarshal([]byte(request), &callbackBody)
	if err != nil {
		return
	}
	a.Token = callbackBody.AccessToken
	a.Refresh = callbackBody.RefreshToken
	a.ExpiresAt = time.Now().Add(time.Second * time.Duration(callbackBody.ExpiresIn))
}

func (a *SpotifyApi) refreshAccess() {
	values := url.Values{}

	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", a.Refresh)
	request, err := a.basicRequest("https://accounts.spotify.com/api/token", values)
	if err != nil {
		return
	}

	refresh := SpotifyRefresh{}
	err = json.Unmarshal([]byte(request), &refresh)
	if err != nil {
		return
	}
	a.Token = refresh.AccessToken
	a.ExpiresAt = time.Now().Add(time.Second * time.Duration(refresh.ExpiresIn))

}

type SpotifyRefresh struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
}

func (a *SpotifyApi) basicRequest(path string, values url.Values) (string, error) {

	id := os.Getenv("client")
	secret := os.Getenv("secret")

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
