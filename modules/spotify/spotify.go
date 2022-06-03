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
	"strconv"
	"strings"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/plugin"
)

var Module Spotify

type Spotify struct {
	plugin.Module
	api      SpotifyApi
	Accounts map[string]SpotifyApi
	id       string
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

func (s *Spotify) PutAttribute(key string) func(str string) error {
	return func(str string) error {
		switch key {
		case "current":
			break
		case "cmd":
			switch str {
			case "next":
				_, err := s.api.next()
				if err != nil {
					return err
				}
			case "previous":
				_, err := s.api.previous()
				if err != nil {
					return err
				}
			}
			break
		case "playing":
			parseBool, err := strconv.ParseBool(str)
			if err != nil {
				return err
			}

			if parseBool {
				_, err = s.api.play()
			} else {
				_, err = s.api.pause()
			}

			if err != nil {
				return err
			}
			err = s.Attributes.Set(s.id, "playing", str)

			break
		}
		return nil
	}
}

func (s *Spotify) GetAttribute(key string) func() (string, error) {
	return func() (string, error) {
		switch key {
		case "current":
			return s.api.current()
		case "cmd":
			return "none", nil
		case "playing":
			song, err := s.api.current()
			if err != nil {
				return "false", err
			}
			c := CurrentlyPlaying{}
			err = json.Unmarshal([]byte(song), &c)
			if err != nil {
				return "false", err
			}
			res := "false"
			if c.Playing {
				res = "true"
			}
			return res, nil
		}
		return "", nil
	}
}

type SpotifyState struct {
	Title      string    `json:"title"`
	Cover      string    `json:"cover"`
	Thumbnail  string    `json:"thumbnail"`
	Artists    string    `json:"artists"`
	Album      string    `json:"album"`
	Progress   int       `json:"progress"`
	Updated    time.Time `json:"updated"`
	Duration   int       `json:"duration"`
	Explicit   bool      `json:"explicit"`
	Playing    bool      `json:"playing"`
	Popularity int       `json:"popularity"`
	Volume     int       `json:"volume"`
	Device     string    `json:"device"`
}

func (s *Spotify) Setup() (plugin.Config, error) {

	s.Frequency = 5000 * time.Millisecond
	return s.Config, nil
}

type CurrentlyPlaying struct {
	Playing bool `json:"is_playing"`
}
type CurrentResponse struct {
	Device struct {
		Id               string `json:"id"`
		IsActive         bool   `json:"is_active"`
		IsPrivateSession bool   `json:"is_private_session"`
		IsRestricted     bool   `json:"is_restricted"`
		Name             string `json:"name"`
		Type             string `json:"type"`
		VolumePercent    int    `json:"volume_percent"`
	} `json:"device"`
	ShuffleState bool        `json:"shuffle_state"`
	RepeatState  string      `json:"repeat_state"`
	Timestamp    int64       `json:"timestamp"`
	Context      interface{} `json:"context"`
	ProgressMs   int         `json:"progress_ms"`
	Item         struct {
		Album struct {
			AlbumType string `json:"album_type"`
			Artists   []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				Id   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				Uri  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			ExternalUrls     struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			Id     string `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				Url    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name                 string `json:"name"`
			ReleaseDate          string `json:"release_date"`
			ReleaseDatePrecision string `json:"release_date_precision"`
			TotalTracks          int    `json:"total_tracks"`
			Type                 string `json:"type"`
			Uri                  string `json:"uri"`
		} `json:"album"`
		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			Id   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			Uri  string `json:"uri"`
		} `json:"artists"`
		AvailableMarkets []string `json:"available_markets"`
		DiscNumber       int      `json:"disc_number"`
		DurationMs       int      `json:"duration_ms"`
		Explicit         bool     `json:"explicit"`
		ExternalIds      struct {
			Isrc string `json:"isrc"`
		} `json:"external_ids"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href        string `json:"href"`
		Id          string `json:"id"`
		IsLocal     bool   `json:"is_local"`
		Name        string `json:"name"`
		Popularity  int    `json:"popularity"`
		PreviewUrl  string `json:"preview_url"`
		TrackNumber int    `json:"track_number"`
		Type        string `json:"type"`
		Uri         string `json:"uri"`
	} `json:"item"`
	CurrentlyPlayingType string `json:"currently_playing_type"`
	Actions              struct {
		Disallows struct {
			Resuming              bool `json:"resuming"`
			TogglingRepeatContext bool `json:"toggling_repeat_context"`
			TogglingRepeatTrack   bool `json:"toggling_repeat_track"`
			TogglingShuffle       bool `json:"toggling_shuffle"`
		} `json:"disallows"`
	} `json:"actions"`
	IsPlaying bool `json:"is_playing"`
}

func (s *Spotify) Update() error {
	if time.Since(s.Module.LastUpdate) >= s.Frequency {
		s.Module.LastUpdate = time.Now()
		return s.push()
	}
	return nil
}

func (s *Spotify) push() error {

	song, err := s.api.current()
	if err != nil {
		return err
	}

	if !json.Valid([]byte(song)) {
		err = s.Attributes.Set(s.id, "playing", "false")
		if err != nil {
			return err
		}
		return nil
	}

	c := CurrentResponse{}
	err = json.Unmarshal([]byte(song), &c)
	if err != nil {
		log.Err(err)
		return err
	}

	sp := SpotifyState{}
	sp.Duration = c.Item.DurationMs
	sp.Progress = c.ProgressMs
	sp.Updated = time.Now()
	sp.Playing = c.IsPlaying
	sp.Album = c.Item.Album.Name
	sp.Title = c.Item.Name
	sp.Popularity = c.Item.Popularity
	sp.Explicit = c.Item.Explicit
	sp.Volume = c.Device.VolumePercent
	sp.Device = c.Device.Name

	var names []string
	for _, artist := range c.Item.Album.Artists {
		names = append(names, artist.Name)
	}

	sp.Artists = strings.Join(names, ", ")

	for _, image := range c.Item.Album.Images {
		switch image.Width {
		case 640:
			sp.Cover = image.Url
		case 64:
			sp.Thumbnail = image.Url
		default:
			continue
		}
	}

	marshal, err := json.Marshal(sp)
	if err != nil {
		return err
	}
	if s.id == "" {
		return nil
	}
	err = s.Attributes.Set(s.id, "current", string(marshal))
	if err != nil {
		return err
	}
	res := "false"
	if sp.Playing {
		s.Frequency = time.Millisecond * 3000
		res = "true"
	} else {
		s.Frequency = time.Millisecond * 15000
	}
	err = s.Attributes.Set(s.id, "playing", res)
	if err != nil {
		return err
	}
	return nil
}

func (s *Spotify) Run() error {

	e := &domain.Entity{
		Name:   "remote",
		Type:   "media",
		Module: "spotify",
	}
	err := s.Entities.Register(e)
	if err != nil {
		return err
	}

	current := &domain.Attribute{
		Key:     "current",
		Value:   "{}",
		Request: "{}",
		Type:    "media",
		Entity:  e.Id,
		Channel: make(chan domain.Attribute),
	}
	go func() {
		for attribute := range current.Channel {
			err := s.PutAttribute(current.Key)(attribute.Request)
			if err != nil {
				log.Err(err)
				return
			}
		}
	}()
	err = s.Attributes.Register(current)
	if err != nil {
		return err
	}

	playing := &domain.Attribute{
		Key:     "playing",
		Value:   "false",
		Request: "false",
		Entity:  e.Id,
		Channel: make(chan domain.Attribute),
	}
	go func() {
		for attribute := range playing.Channel {
			err := s.PutAttribute(playing.Key)(attribute.Request)
			if err != nil {
				log.Err(err)
				return
			}
		}
	}()

	err = s.Attributes.Register(playing)
	if err != nil {
		return err
	}

	cmd := &domain.Attribute{
		Key:     "cmd",
		Value:   "none",
		Request: "none",
		Entity:  e.Id,
		Channel: make(chan domain.Attribute),
	}
	go func() {
		for attribute := range cmd.Channel {
			err := s.PutAttribute(cmd.Key)(attribute.Request)
			if err != nil {
				log.Err(err)
				return
			}
		}
	}()

	err = s.Attributes.Register(cmd)
	if err != nil {
		return err
	}

	if e.Id != "" {
		s.id = e.Id
		if e.Config == "" {
			a := SpotifyApi{}
			m, _ := json.Marshal(&a)
			e.Config = string(m)
			s.api.Authenticate()
			err = s.Entities.Config(s.id, string(m))
			if err != nil {
				return err
			}
			return nil
		}
		a := SpotifyApi{}
		err = json.Unmarshal([]byte(e.Config), &a)
		if err != nil {
			return err
		}
		s.api = a

	}

	s.api.Authenticate()
	marshal, err := json.Marshal(s.api)
	if err != nil {
		return err
	}
	err = s.Entities.Config(s.id, string(marshal))
	if err != nil {
		return err
	}

	return nil
}

func (s *SpotifyApi) current() (string, error) {
	return s.authenticatedRequest("GET", "me/player")
}

func (s *SpotifyApi) play() (string, error) {
	return s.authenticatedRequest("PUT", "me/player/play")
}

func (s *SpotifyApi) pause() (string, error) {
	return s.authenticatedRequest("PUT", "me/player/pause")
}

func (s *SpotifyApi) next() (string, error) {
	return s.authenticatedRequest("POST", "me/player/next")
}

func (s *SpotifyApi) previous() (string, error) {
	return s.authenticatedRequest("POST", "me/player/previous")
}

func (s *SpotifyApi) analyzeTrack(trackId string) (string, error) {
	return s.authenticatedRequest("POST", fmt.Sprintf("audio-analysis/%s", trackId))
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

	scope := "user-modify-playback-state user-read-currently-playing user-read-playback-state"
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
