// Copyright (c) 2021 Braden Nicholson

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
)

var Module Weather

type Weather struct {
	plugin.Module
	forecast WeatherAPI
	eId      string
}

const weatherUrl = "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&hourly=temperature_2m,relativehumidity_2m,precipitation,weathercode&daily=weathercode,temperature_2m_max,temperature_2m_min,sunrise,sunset,precipitation_sum&current_weather=true&temperature_unit=fahrenheit&windspeed_unit=mph&precipitation_unit=inch&timeformat=unixtime&timezone=America%%2FLos_Angeles"

type WeatherAPI struct {
	UtcOffsetSeconds int     `json:"utc_offset_seconds"`
	GenerationtimeMs float64 `json:"generationtime_ms"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Elevation        float64 `json:"elevation"`
	CurrentWeather   struct {
		Temperature   float64 `json:"temperature"`
		Winddirection float64 `json:"winddirection"`
		Weathercode   float64 `json:"weathercode"`
		Time          float64 `json:"time"`
		Windspeed     float64 `json:"windspeed"`
	} `json:"current_weather"`
	Hourly      Hourly      `json:"hourly"`
	Daily       Daily       `json:"daily"`
	HourlyUnits HourlyUnits `json:"hourly_units"`
	DailyUnits  DailyUnits  `json:"daily_units"`
}

type CurrentWeather struct {
	Temperature     float64 `json:"temperature"`
	TemperatureNext float64 `json:"temperatureNext"`
	Humidity        float64 `json:"humidity"`
	Min             float64 `json:"min"`
	Max             float64 `json:"max"`
}

type Daily struct {
	Sunrise          []float64 `json:"sunrise"`
	PrecipitationSum []float64 `json:"precipitation_sum"`
	Weathercode      []float64 `json:"weathercode"`
	Temperature2MMin []float64 `json:"temperature_2m_min"`
	Time             []float64 `json:"time"`
	Temperature2MMax []float64 `json:"temperature_2m_max"`
	Sunset           []float64 `json:"sunset"`
}

type Hourly struct {
	ShortwaveRadiation  []float64 `json:"shortwave_radiation"`
	Precipitation       []float64 `json:"precipitation"`
	Relativehumidity2M  []float64 `json:"relativehumidity_2m"`
	Winddirection10M    []float64 `json:"winddirection_10m"`
	Weathercode         []float64 `json:"weathercode"`
	Windgusts10M        []float64 `json:"windgusts_10m"`
	ApparentTemperature []float64 `json:"apparent_temperature"`
	Time                []float64 `json:"time"`
	Windspeed10M        []float64 `json:"windspeed_10m"`
	Temperature2M       []float64 `json:"temperature_2m"`
}

type HourlyUnits struct {
	ShortwaveRadiation  string `json:"shortwave_radiation"`
	Precipitation       string `json:"precipitation"`
	Winddirection10M    string `json:"winddirection_10m"`
	Windspeed10M        string `json:"windspeed_10m"`
	ApparentTemperature string `json:"apparent_temperature"`
	Weathercode         string `json:"weathercode"`
	Windgusts10M        string `json:"windgusts_10m"`
	Time                string `json:"time"`
	Temperature2M       string `json:"temperature_2m"`
	Relativehumidity2M  string `json:"relativehumidity_2m"`
}

type DailyUnits struct {
	Sunrise          string `json:"sunrise"`
	PrecipitationSum string `json:"precipitation_sum"`
	Weathercode      string `json:"weathercode"`
	Temperature2MMin string `json:"temperature_2m_min"`
	Time             string `json:"time"`
	Temperature2MMax string `json:"temperature_2m_max"`
	Sunset           string `json:"sunset"`
}

func init() {
	config := plugin.Config{
		Name:        "weather",
		Type:        "module",
		Description: "Live weather updates",
		Version:     "0.2.1",
		Author:      "Braden Nicholson",
	}
	Module.forecast = WeatherAPI{}
	Module.eId = ""
	Module.Config = config
}
func (v *Weather) currentBuffer() (string, error) {

	index := time.Now().Hour()

	if len(v.forecast.Hourly.Temperature2M) < index {
		cw := CurrentWeather{}
		marshal, err := json.Marshal(cw)
		if err != nil {
			return "", err
		}
		return string(marshal), nil
	}

	cw := CurrentWeather{
		Temperature:     v.forecast.Hourly.Temperature2M[index],
		TemperatureNext: v.forecast.Hourly.Temperature2M[(index+1)%24],
		Humidity:        v.forecast.Hourly.Relativehumidity2M[index],
		Min:             v.forecast.Daily.Temperature2MMin[0],
		Max:             v.forecast.Daily.Temperature2MMax[0],
	}

	marshal, err := json.Marshal(cw)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

func (v *Weather) forecastBuffer() (string, error) {
	marshal, err := json.Marshal(v.forecast)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

func (v *Weather) fetchWeather() error {

	lat, flat := os.LookupEnv("weatherLat")
	if !flat {
		return fmt.Errorf("WEATHER: env weatherLat not set")
	}
	lon, flon := os.LookupEnv("weatherLon")
	if !flon {
		return fmt.Errorf("WEATHER: env weatherLat not set")
	}

	request, err := http.NewRequest("GET", fmt.Sprintf(weatherUrl, lat, lon), nil)
	if err != nil {
		return err
	}

	cli := http.Client{}
	cli.Timeout = time.Second * 5
	defer cli.CloseIdleConnections()

	do, err := cli.Do(request)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(do.Body)
	if err != nil {
		return err
	}

	w := WeatherAPI{}
	err = json.Unmarshal(buf.Bytes(), &w)
	if err != nil {
		return err
	}

	v.forecast = w

	_ = do.Body.Close()

	return nil

}

func (v *Weather) Setup() (plugin.Config, error) {
	err := v.UpdateInterval(15000)
	if err != nil {
		return plugin.Config{}, err
	}
	return v.Config, nil
}

func (v *Weather) pull() error {
	err := v.fetchWeather()
	if err != nil {
		return err
	}
	buffer, err := v.forecastBuffer()
	if err != nil {
		return err
	}
	err = v.Attributes.Set(v.eId, "forecast", buffer)
	if err != nil {
		return err
	}

	current, err := v.currentBuffer()
	if err != nil {
		return err
	}
	err = v.Attributes.Set(v.eId, "weather", current)
	if err != nil {
		return err
	}
	return nil
}
func (v *Weather) Update() error {
	if v.Ready() {
		err := v.pull()
		if err != nil {
			v.ErrF("Weather fetch failed: %s", err.Error())
			return err
		}
	}
	return nil
}

func (v *Weather) Run() error {

	e := &domain.Entity{
		Name:   "weather",
		Module: "weather",
		Type:   "media",
	}
	err := v.Entities.Register(e)
	if err != nil {
		return err
	}

	buffer, err := v.forecastBuffer()
	if err != nil {
		return err
	}

	currentBuffer, err := v.currentBuffer()
	if err != nil {
		return err
	}

	current := &domain.Attribute{
		Key:     "weather",
		Value:   currentBuffer,
		Request: currentBuffer,
		Type:    "media",
		Order:   0,
		Entity:  e.Id,
	}

	forecast := &domain.Attribute{
		Key:     "forecast",
		Value:   buffer,
		Request: buffer,
		Type:    "media",
		Order:   0,
		Entity:  e.Id,
	}
	v.eId = e.Id

	err = v.Attributes.Register(forecast)
	if err != nil {
		return err
	}
	err = v.Attributes.Register(current)
	if err != nil {
		return err
	}

	return nil
}
