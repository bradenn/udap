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

var Export Spotify

type SpotifyEnv struct {
	Client string `json:"client"`
	Secret string `json:"secret"`
}

type Spotify struct {
}

func (s Spotify) InitInstance() (string, error) {
	spotify := SpotifyApi{}
	spotify.Authenticate()

	marshal, err := json.Marshal(spotify)
	if err != nil {
		return "", err
	}

	return string(marshal), nil
}

func (s Spotify) Initialize(env string) {
	spotifyEnv := SpotifyEnv{}
	err := json.Unmarshal([]byte(env), &spotifyEnv)
	if err != nil {
		return
	}
	err = os.Setenv("client", spotifyEnv.Client)
	if err != nil {
		return
	}
	err = os.Setenv("secret", spotifyEnv.Secret)
	if err != nil {
		return
	}
}

func (s Spotify) Metadata() template.Metadata {
	metadata := template.Metadata{
		Name:        "Spotify",
		Description: "Remote control for spotify",
		Version:     "1.0.1",
		Author:      "Braden Nicholson",
	}
	return metadata
}

func (s Spotify) Poll(v string) (string, error) {
	spotifyApi := SpotifyApi{}
	err := json.Unmarshal([]byte(v), &spotifyApi)
	if err != nil {
		return "", err
	}
	return spotifyApi.CurrentSong()
}

func (s Spotify) Run(v string, action string) (string, error) {
	spotifyApi := SpotifyApi{}
	err := json.Unmarshal([]byte(v), &spotifyApi)
	if err != nil {
		return "", err
	}
	switch action {
	case "play":
		return spotifyApi.Play()
	case "pause":
		return spotifyApi.Pause()
	}
	return "", nil
}

func (s *SpotifyApi) CurrentSong() (string, error) {
	return s.authenticatedRequest("GET", "me/player/currently-playing")
}

func (s *SpotifyApi) Play() (string, error) {
	return s.authenticatedRequest("PUT", "me/player/play")
}

func (s *SpotifyApi) Pause() (string, error) {
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
}

func (s *SpotifyApi) Authenticate() {
	if s.Token == "" {
		fmt.Printf("Please log into your spotify account here: %s", s.loginURL())

		e := make(chan string)
		go s.beginListening(e)
		msg := <-e
		s.requestToken(msg)
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

func (s *SpotifyApi) loginURL() string {

	scope := "user-modify-playback-state user-read-currently-playing"
	clientId := os.Getenv("client")

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
