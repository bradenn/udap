// Copyright (c) 2023 Braden Nicholson

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/plugin"
)

var Module Google

const (
	CommandSetHeat  = "sdm.devices.commands.ThermostatTemperatureSetpoint.SetHeat"
	CommandSetCool  = "sdm.devices.commands.ThermostatTemperatureSetpoint.SetCool"
	CommandSetRange = "sdm.devices.commands.ThermostatTemperatureSetpoint.SetRange"
	CommandSetMode  = "sdm.devices.commands.ThermostatMode.SetMode"
)

type Google struct {
	plugin.Module
	request  chan domain.Attribute
	oauth    sync.Mutex
	entityId string
	deviceId string

	mode string
}

type SetMode struct {
	Heat float64 `json:"heat,omitempty"`
	Cool float64 `json:"cool,omitempty"`
}

type ModeSet struct {
	Mode string `json:"mode,omitempty"`
}

func init() {

	config := plugin.Config{
		Name:        "googleHome",
		Type:        "module",
		Description: "Google Smart Home products",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func (c *Google) mux() {
	for attribute := range c.request {
		c.handleRequest(attribute)
	}
}

func convertCToF(c float64) float64 {
	return c*(9.0/5.0) + 32.0
}

func convertFToC(f float64) float64 {
	return (f - 32.0) * (5.0 / 9.0)
}

type HeatCoolParams struct {
	HeatCelsius float64 `json:"heatCelsius,omitempty"`
	CoolCelsius float64 `json:"coolCelsius,omitempty"`
}

type Command struct {
	Command string          `json:"command"`
	Params  json.RawMessage `json:"params"`
}

func (c *Command) SetHeat(heat float64) {
	hcp := HeatCoolParams{
		HeatCelsius: convertFToC(heat),
	}
	marshal, err := json.Marshal(hcp)
	if err != nil {
		return
	}
	c.Command = CommandSetHeat
	c.Params = marshal
}

func (c *Command) SetCool(cool float64) {
	hcp := HeatCoolParams{
		CoolCelsius: convertFToC(cool),
	}
	marshal, err := json.Marshal(hcp)
	if err != nil {
		return
	}
	c.Command = CommandSetCool
	c.Params = marshal
}

func (c *Command) SetMode(mode string) {
	hcp := ModeSet{
		Mode: mode,
	}
	if mode != "HEAT" && mode != "COOL" && mode != "OFF" && mode != "HEATCOOL" {
		hcp = ModeSet{
			Mode: "OFF",
		}
	}
	marshal, err := json.Marshal(hcp)
	if err != nil {
		return
	}
	c.Command = CommandSetMode
	c.Params = marshal
}

func (c *Google) setBestGuess(attribute domain.Attribute) error {
	mode, err := c.Attributes.FindByComposite(attribute.Entity, "mode")
	if err != nil {
		return err
	}
	if mode.Value == "HEAT" {
		return c.setHeat(attribute)
	} else if mode.Value == "COOL" {
		return c.setCool(attribute)
	} else {

	}

	return nil
}

func (c *Google) setHeat(attribute domain.Attribute) error {
	float, err := strconv.ParseFloat(attribute.Request, 64)
	if err != nil {
		return err
	}

	cmd := Command{}
	cmd.SetHeat(float)

	marshal, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	_, err = c.sendAuthenticatedRequest(marshal)
	if err != nil {
		return err
	}
	return nil
	//fmt.Println(string(request))
}
func (c *Google) setMode(attribute domain.Attribute) {
	cmd := Command{}
	cmd.SetMode(attribute.Request)

	marshal, err := json.Marshal(cmd)
	if err != nil {
		log.Err(err)
		return
	}

	_, err = c.sendAuthenticatedRequest(marshal)
	if err != nil {
		log.Err(err)
		return
	}

	err = c.forceUpdate()
	if err != nil {
		return
	}

	//fmt.Println(string(request))

}
func (c *Google) setCool(attribute domain.Attribute) error {
	float, err := strconv.ParseFloat(attribute.Request, 64)
	if err != nil {
		return err
	}

	cmd := Command{}
	cmd.SetCool(float)

	marshal, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	_, err = c.sendAuthenticatedRequest(marshal)
	if err != nil {
		return err
	}
	return nil
	//fmt.Println(string(request))

}

func (c *Google) handleRequest(attribute domain.Attribute) {
	switch attribute.Key {
	case "target":
		err := c.setBestGuess(attribute)
		if err != nil {
			fmt.Println(err)
			break
		}
		err = c.Attributes.Set(c.entityId, "target", attribute.Request)
		if err != nil {
			log.Err(err)
		}
		break
	case "heat":
		c.setHeat(attribute)
		break
	case "mode":
		c.setMode(attribute)
		err := c.Attributes.Set(c.entityId, "mode", attribute.Request)
		if err != nil {
			log.Err(err)
		}
		break
	case "cool":
		c.setCool(attribute)
		break
	default:
		break
	}
}

const refreshURL = "https://www.googleapis.com/oauth2/v4/token?client_id=%s&client_secret=%s&refresh_token=%s&grant_type=refresh_token"
const commandURL = "https://smartdevicemanagement.googleapis.com/v1/%s:executeCommand"
const apiURL = "https://smartdevicemanagement.googleapis.com/v1/enterprises/%s/devices"
const authTokenURL = "https://www.googleapis.com/oauth2/v4/token?client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&redirect_uri=%s"
const pcmURL = "https://nestservices.google.com/partnerconnections/%s/auth?redirect_uri=%s&access_type=offline&prompt=consent&client_id=%s&response_type=code&scope=https://www.googleapis.com/auth/sdm.service"
const redirectUrl = "https://google-oauth.udap.app"

func (c *Google) Setup() (plugin.Config, error) {
	err := c.UpdateInterval(30000)
	if err != nil {
		return plugin.Config{}, err
	}
	return Module.Config, nil
}

func generateServicePartnerURL() (string, error) {

	project, foundProjectId := os.LookupEnv("googleProjectId")
	if !foundProjectId {
		return "", fmt.Errorf("project id environment variable not set")
	}
	client, foundClientId := os.LookupEnv("googleClientId")
	if !foundClientId {
		return "", fmt.Errorf("client id environment variable not set")
	}
	return fmt.Sprintf(pcmURL, project, redirectUrl, client), nil
}

func generateTokenURL(code string) (string, error) {

	client, foundClientId := os.LookupEnv("googleClientId")
	if !foundClientId {
		return "", fmt.Errorf("client id environment variable not set")
	}

	secret, foundSecret := os.LookupEnv("googleClientSecret")
	if !foundSecret {
		return "", fmt.Errorf("client id environment variable not set")
	}

	return fmt.Sprintf(authTokenURL, client, secret, code, redirectUrl), nil
}

func (c *Google) runOAuth() error {
	if !c.oauth.TryLock() {
		return fmt.Errorf("oauth is already in progress. Please wait")
	}
	defer c.oauth.Unlock()

	authCode := make(chan string, 1)

	pUrl, err := generateServicePartnerURL()
	if err != nil {
		log.Err(err)
		return err
	}

	log.Critical("Google Home Login: %s", pUrl)

	srv := http.Server{}
	router := chi.NewRouter()
	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {

		code := request.URL.Query().Get("code")

		//fmt.Printf("Got code: %s\n", code)

		_, err = writer.Write([]byte("You have been signed in, please come again :)"))
		if err != nil {
			return
		}

		authCode <- code
	})
	srv.Handler = router

	go func() {
		srv.Addr = ":8976"

		err = srv.ListenAndServe()
		if err != nil {
			log.Err(err)
		}
		fmt.Println("Server closed.")
	}()

	// Wait for the auth code
	auth := <-authCode
	// Shutdown the server
	srv.Shutdown(context.Background())
	// Store the auth code
	err = c.SetConfig("auth", auth)
	if err != nil {
		return err
	}

	return nil
}

func (c *Google) fetchToken() error {
	// Get the auth code from storage
	config, err := c.GetConfig("auth")
	// Generate the token fetching url
	url, err := generateTokenURL(config)
	if err != nil {
		return err
	}
	// Print out that URL for reference
	fmt.Println(url)

	cli := http.Client{}
	// Send the request
	post, err := cli.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	// Read the stream to an array of bytes
	all, err := io.ReadAll(post.Body)
	if err != nil {
		return err
	}
	// Close the stream
	defer post.Body.Close()

	// Prepare the response object
	rt := TokenResponse{}
	// Read the response into the struct
	err = json.Unmarshal(all, &rt)
	if err != nil {
		return err
	}

	// Store the new values

	err = c.SetConfig("token", rt.AccessToken)
	if err != nil {
		return err
	}

	err = c.SetConfig("refresh", rt.RefreshToken)
	if err != nil {
		return err
	}

	exp := fmt.Sprintf("%d", time.Now().Add(time.Duration(rt.ExpiresIn)*time.Second).Unix())

	err = c.SetConfig("expires", exp)
	if err != nil {
		return err
	}

	return nil

}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func (c *Google) shouldRefresh() bool {
	refresh, err := c.GetConfig("refresh")
	if err != nil {
		return false
	}

	expires, err := c.GetConfig("expires")
	if err != nil {
		return false
	}

	parseInt, err := strconv.ParseInt(expires, 10, 64)
	if err != nil {
		return false
	}

	return (refresh != "unset") && time.Now().After(time.Unix(parseInt, 0))
}

func (c *Google) sendAuthenticatedRequest(data []byte) ([]byte, error) {
	if c.shouldRefresh() {
		err := c.refresh()
		if err != nil {
			return []byte{}, err
		}
	}

	var buf bytes.Buffer
	buf.Write(data)

	request, err := http.NewRequest("POST", fmt.Sprintf(commandURL, c.deviceId), &buf)
	if err != nil {
		return []byte{}, err
	}

	token, err := c.GetConfig("token")
	if err != nil {
		log.Err(err)
		return []byte{}, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := http.Client{}
	defer client.CloseIdleConnections()
	response, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}

	var res bytes.Buffer

	_, err = res.ReadFrom(response.Body)
	if err != nil {
		return []byte{}, err
	}
	_ = response.Body.Close()
	return res.Bytes(), nil
}

func (c *Google) sendAPIRequest() ([]byte, error) {
	if c.shouldRefresh() {
		err := c.refresh()
		if err != nil {
			return []byte{}, err
		}
	}

	var buf bytes.Buffer
	buf.Write([]byte(""))
	project, foundProjectId := os.LookupEnv("googleProjectId")
	if !foundProjectId {
		return []byte{}, fmt.Errorf("project id environment variable not set")
	}

	request, err := http.NewRequest("GET", fmt.Sprintf(apiURL, project), nil)
	if err != nil {
		return []byte{}, err
	}

	token, err := c.GetConfig("token")
	if err != nil {
		log.Err(err)
		return []byte{}, err
	}

	request.Header.Set("Content-Type", "application/json")
	//fmt.Println(token)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := http.Client{}
	defer client.CloseIdleConnections()
	response, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}

	var res bytes.Buffer

	_, err = res.ReadFrom(response.Body)
	if err != nil {
		return []byte{}, err
	}
	//fmt.Println(res.String())
	_ = response.Body.Close()
	return res.Bytes(), nil
}

type Device struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Assignee string `json:"assignee"`
	Traits   struct {
		SdmDevicesTraitsInfo struct {
			CustomName string `json:"customName"`
		} `json:"sdm.devices.traits.Info"`
		SdmDevicesTraitsHumidity struct {
			AmbientHumidityPercent float64 `json:"ambientHumidityPercent"`
		} `json:"sdm.devices.traits.Humidity"`
		SdmDevicesTraitsConnectivity struct {
			Status string `json:"status"`
		} `json:"sdm.devices.traits.Connectivity"`
		SdmDevicesTraitsFan struct {
			TimerMode string `json:"timerMode"`
		} `json:"sdm.devices.traits.Fan"`
		SdmDevicesTraitsThermostatMode struct {
			Mode           string   `json:"mode"`
			AvailableModes []string `json:"availableModes"`
		} `json:"sdm.devices.traits.ThermostatMode"`
		SdmDevicesTraitsThermostatEco struct {
			AvailableModes []string `json:"availableModes"`
			Mode           string   `json:"mode"`
			HeatCelsius    float64  `json:"heatCelsius"`
			CoolCelsius    float64  `json:"coolCelsius"`
		} `json:"sdm.devices.traits.ThermostatEco"`
		SdmDevicesTraitsThermostatHvac struct {
			Status string `json:"status"`
		} `json:"sdm.devices.traits.ThermostatHvac"`
		SdmDevicesTraitsSettings struct {
			TemperatureScale string `json:"temperatureScale"`
		} `json:"sdm.devices.traits.Settings"`
		SdmDevicesTraitsThermostatTemperatureSetpoint struct {
			CoolCelsius float64 `json:"coolCelsius"`
			HeatCelsius float64 `json:"heatCelsius"`
		} `json:"sdm.devices.traits.ThermostatTemperatureSetpoint"`
		SdmDevicesTraitsTemperature struct {
			AmbientTemperatureCelsius float64 `json:"ambientTemperatureCelsius"`
		} `json:"sdm.devices.traits.Temperature"`
	} `json:"traits"`
	ParentRelations []struct {
		Parent      string `json:"parent"`
		DisplayName string `json:"displayName"`
	} `json:"parentRelations"`
}

type QueryResponse struct {
	Devices []Device `json:"devices"`
}

func (c *Google) forceUpdate() error {
	dat, err := c.sendAPIRequest()
	if err != nil {
		return err
	}

	qr := QueryResponse{}
	err = json.Unmarshal(dat, &qr)
	if err != nil {
		return err
	}

	err = c.updateValues(qr)
	if err != nil {
		return err
	}

	return nil
}

func (c *Google) Update() error {
	if !c.Ready() {
		return nil
	}
	dat, err := c.sendAPIRequest()
	if err != nil {
		return err
	}

	qr := QueryResponse{}
	err = json.Unmarshal(dat, &qr)
	if err != nil {
		return err
	}

	err = c.updateValues(qr)
	if err != nil {
		return err
	}

	return nil
}

type Thermostat struct {
	Connected   bool    `json:"connected"`
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
	Mode        string  `json:"mode"`
	EcoCool     float64 `json:"ecoCool"`
	EcoHeat     float64 `json:"ecoHeat"`
}

type ThermostatSensor struct {
	Humidity float64 `json:"humidity"`
	Temp     float64 `json:"temp"`
}

func (c *Google) updateValues(qr QueryResponse) error {

	for _, device := range qr.Devices {
		if device.Type != "sdm.devices.types.THERMOSTAT" {
			continue
		}
		connected := false
		if device.Traits.SdmDevicesTraitsConnectivity.Status == "CONNECTED" {
			connected = true
		}
		ts := Thermostat{
			Humidity:    device.Traits.SdmDevicesTraitsHumidity.AmbientHumidityPercent,
			Temperature: convertCToF(device.Traits.SdmDevicesTraitsTemperature.AmbientTemperatureCelsius),
			Mode:        device.Traits.SdmDevicesTraitsThermostatMode.Mode,
			EcoCool:     convertCToF(device.Traits.SdmDevicesTraitsThermostatEco.CoolCelsius),
			EcoHeat:     convertCToF(device.Traits.SdmDevicesTraitsThermostatEco.HeatCelsius),
			Connected:   connected,
		}

		marshal, err := json.Marshal(ts)
		if err != nil {
			return err
		}

		err = c.Attributes.Set(c.entityId, "thermostat", string(marshal))
		if err != nil {
			return err
		}

		err = c.Attributes.Set(c.entityId, "heat", fmt.Sprintf("%f", convertCToF(device.Traits.SdmDevicesTraitsThermostatTemperatureSetpoint.HeatCelsius)))
		if err != nil {
			return err
		}

		sensorData := ThermostatSensor{
			Temp:     device.Traits.SdmDevicesTraitsTemperature.AmbientTemperatureCelsius,
			Humidity: ts.Humidity,
		}

		sensorJSON, err := json.Marshal(sensorData)
		if err != nil {
			return err
		}

		err = c.Attributes.Set(c.entityId, "sensor", string(sensorJSON))
		if err != nil {
			return err
		}

		err = c.Attributes.Set(c.entityId, "cool", fmt.Sprintf("%f", convertCToF(device.Traits.SdmDevicesTraitsThermostatTemperatureSetpoint.CoolCelsius)))
		if err != nil {
			return err
		}

		err = c.Attributes.Set(c.entityId, "mode", device.Traits.SdmDevicesTraitsThermostatMode.Mode)
		if err != nil {
			return err
		}

		if device.Traits.SdmDevicesTraitsThermostatMode.Mode == "HEAT" {
			err = c.Attributes.Set(c.entityId, "target", fmt.Sprintf("%f", convertCToF(device.Traits.SdmDevicesTraitsThermostatTemperatureSetpoint.HeatCelsius)))
			if err != nil {
				return err
			}
		} else if device.Traits.SdmDevicesTraitsThermostatMode.Mode == "COOL" {
			err = c.Attributes.Set(c.entityId, "target", fmt.Sprintf("%f", convertCToF(device.Traits.SdmDevicesTraitsThermostatTemperatureSetpoint.CoolCelsius)))
			if err != nil {
				return err
			}
		}

	}
	return nil
}
func (c *Google) initDevices() error {
	dat, err := c.sendAPIRequest()
	if err != nil {
		return err
	}

	qr := QueryResponse{}
	err = json.Unmarshal(dat, &qr)
	if err != nil {
		return err
	}

	for _, device := range qr.Devices {
		c.deviceId = device.Name
		if device.Type != "sdm.devices.types.THERMOSTAT" {
			continue
		}
		dev := domain.Entity{
			Name:   strings.ToLower("thermostat"),
			Type:   "thermostat",
			Module: "google",
		}
		err = c.Entities.Register(&dev)
		if err != nil {
			return err
		}
		c.entityId = dev.Id
		channel := c.request
		media := &domain.Attribute{
			Key:     "thermostat",
			Value:   "{}",
			Request: "{}",
			Order:   0,
			Type:    "media",
			Entity:  dev.Id,
			Channel: channel,
		}
		err = c.Attributes.Register(media)
		if err != nil {
			return err
		}

		sensor := &domain.Attribute{
			Key:     "sensor",
			Value:   "{\"temp\": 0, \"humidity\": 0}",
			Request: "{\"temp\": 0, \"humidity\": 0}",
			Order:   0,
			Type:    "media",
			Serial:  dev.Name,
			Entity:  dev.Id,
			Channel: channel,
		}
		err = c.Attributes.Register(sensor)
		if err != nil {
			return err
		}

		target := &domain.Attribute{
			Key:     "target",
			Value:   "50",
			Request: "50",
			Order:   0,
			Type:    "number",
			Entity:  dev.Id,
			Channel: channel,
		}
		err = c.Attributes.Register(target)
		if err != nil {
			return err
		}

		heat := &domain.Attribute{
			Key:     "heat",
			Value:   "0",
			Request: "0",
			Order:   0,
			Type:    "number",
			Entity:  dev.Id,
			Channel: channel,
		}
		err = c.Attributes.Register(heat)
		if err != nil {
			return err
		}

		cool := &domain.Attribute{
			Key:     "cool",
			Value:   "0",
			Request: "0",
			Order:   0,
			Type:    "number",
			Entity:  dev.Id,
			Channel: channel,
		}
		err = c.Attributes.Register(cool)
		if err != nil {
			return err
		}

		mode := &domain.Attribute{
			Key:     "mode",
			Value:   "OFF",
			Request: "OFF",
			Order:   0,
			Type:    "string",
			Entity:  dev.Id,
			Channel: channel,
		}
		err = c.Attributes.Register(mode)
		if err != nil {
			return err
		}

	}
	err = c.updateValues(qr)
	if err != nil {
		return err
	}

	return nil
}
func (c *Google) Run() error {

	err := c.InitConfig("auth", "unset")
	if err != nil {
		return err
	}

	err = c.InitConfig("token", "unset")
	if err != nil {
		return err
	}

	err = c.InitConfig("refresh", "unset")
	if err != nil {
		return err
	}

	err = c.InitConfig("expires", "0")
	if err != nil {
		return err
	}

	c.request = make(chan domain.Attribute)
	c.oauth = sync.Mutex{}

	auth, err := c.GetConfig("auth")
	if err != nil {
		return err
	}

	token, err := c.GetConfig("token")
	if err != nil {
		return err
	}

	go c.mux()

	if auth == "unset" || token == "unset" {
		go func() {
			errf := c.connectAccount()
			if errf != nil {
				log.Err(errf)
			}
		}()
	} else {
		log.Event("Google Auth Code: %s", auth)
		if c.shouldRefresh() {
			err = c.refresh()
			if err != nil {
				return err
			}
		}
		err = c.initDevices()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Google) connectAccount() error {
	err := c.runOAuth()
	if err != nil {
		return err
	}

	// Get the auth code
	auth, err := c.GetConfig("auth")
	if err != nil {
		return err
	}

	// Print out the results
	log.Event("Code: %s", auth)

	// Try to get a token
	err = c.fetchToken()
	if err != nil {
		return err
	}

	err = c.initDevices()
	if err != nil {
		return err
	}

	return nil
}

func (c *Google) Dispose() error {
	close(c.request)
	return nil
}

type RefreshClass struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (c *Google) refresh() error {
	refresh, err := c.GetConfig("refresh")
	if err != nil {
		return err
	}

	clientId, foundClientId := os.LookupEnv("googleClientId")
	if !foundClientId {
		return fmt.Errorf("client id environment variable not set")
	}

	secret, foundSecret := os.LookupEnv("googleClientSecret")
	if !foundSecret {
		return fmt.Errorf("client id environment variable not set")
	}

	request, err := http.NewRequest("POST", fmt.Sprintf(refreshURL, clientId, secret, refresh), nil)
	if err != nil {
		return err
	}

	//fmt.Println(fmt.Sprintf(refreshURL, clientId, secret, refresh))

	client := http.Client{}
	defer client.CloseIdleConnections()

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	all, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	rc := RefreshClass{}

	err = json.Unmarshal(all, &rc)
	if err != nil {
		return err
	}
	err = c.SetConfig("token", rc.AccessToken)
	if err != nil {
		return err
	}
	err = c.SetConfig("expires", fmt.Sprintf("%d", time.Now().Add(time.Duration(rc.ExpiresIn)*time.Second).Unix()))
	if err != nil {
		return err
	}

	return nil

}
