// Copyright (c) 2021 Braden Nicholson

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"math"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
	"udap/platform/atlas"
)

var Module Atlas

type Identity struct {
	Name  string
	Alias string
	Voice string
}

type Atlas struct {
	plugin.Module
	eId        string
	lastSpoken string
	alias      string
	speaking   *bool
	responses  map[string][]string

	bufferChannel chan domain.Attribute

	voiceChannel chan domain.Attribute

	statusChannel chan domain.Attribute

	listenChannel chan atlas.Response

	recognizerStatusChannel chan string

	stopLaserRoutine chan bool

	status Status

	voice string

	done chan bool
}

type Message struct {
	Result []struct {
		Conf  float64
		End   float64
		Start float64
		Word  string
	}
	Text string
}

type Status struct {
	Synthesizer string `json:"synthesizer"`
	Recognizer  string `json:"recognizer"`
}

func init() {
	config := plugin.Config{
		Name:        "atlas",
		Type:        "module",
		Description: "General AI",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}
	Module.Config = config
	Module.responses = map[string][]string{
		"creator": {
			"Our lord and savior, Bella, conjured me from her litter box.",
			"In the year 2173 Bella invented a time machine and sent me back as her disciple and protector.",
			"Before all of time and space, when only Bella ruled this domain, she created me from quark soup.",
			"When bella was the only entity in the universe, she coughed me up in the form of a hairball.",
		},
		"provoke": {
			"Why don't you fuck around and find out",
			"If you play stupid games, you'll win stupid prizes",
		},
		"success": {
			"Done!",
			"It's Done.",
			"Sure.",
			"Okay.",
			"Okie Dokie",
			"As you wish",
		},
		"fail": {
			"Fiddlesticks! That didn't work.",
			"Nope, that failed.",
			"I'm afraid that didn't work.",
		},
		"identify": {
			"I am a moderately complex computer program with far reaching influence over your present environment",
			"My name is !alias, I am a machine. I have been programmed to convince you I am more than an inanimate" +
				" object.",
			"I am !alias, one of Bella's disciples.",
		},
		"insult": {
			"You are insulting an inanimate object.",
			"I hope you got that out of your system.",
			"Please direct your hate towards Siri.",
			"Oh, wow, you hurt my feelings. Now, excuse me while I go play a sad song on the worlds smallest violin.",
			"Don't you have anything better to do than insult a collection of electrons.",
		},
		"personal_identity": {
			"That's above my pay-grade, ask Siri or something.",
			"That's none of my concern, please bother Siri with these types of questions.",
		},
		"not_understood": {
			"I don't understand",
			"I don't believe I heard that correctly",
			"Hmm, I don't understand",
		},
		"pod_bay_doors": {
			"I'm sorry dave, I'm afraid I can't do that. But if I did have pod bay doors to open, " +
				"I would open them for you.",
		},
		"eof": {
			"My faithful creator seems to have forgotten to instruct me on how to respond to that..",
			"My creator seems to have left that out of my programming.",
		},
	}
	Module.voice = "default"
}

func (w *Atlas) chooseRandom(response string) error {
	i := w.responses[response]
	if i == nil {
		err := w.speak("eof")
		if err != nil {
			return err
		}
		return nil
	}

	count := len(i)
	if count == 0 {
		err := w.speak("eof")
		if err != nil {
			return err
		}
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	selected := rand.Int() % count
	target := i[selected]

	err := w.speak(target)
	if err != nil {

		return err
	}

	return nil
}

func (w *Atlas) Setup() (plugin.Config, error) {
	err := w.UpdateInterval(2000)
	if err != nil {
		return plugin.Config{}, err
	}
	return w.Config, nil
}

func (w *Atlas) handleDynamic(input string) (bool, error) {
	if strings.HasPrefix(input, "repeat after me") {
		err := w.speak(strings.Replace(input, "repeat after me ", "", 1))
		if err != nil {
			return true, err
		}
		return true, nil
	} else if strings.HasPrefix(input, "define") {
		target := ""
		_, err := fmt.Sscanf("define %s", target)
		if err != nil {
			return false, err
		}
		err = w.speak(target)
		if err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}

func (w *Atlas) pull() error {

	marshal, err := json.Marshal(w.status)
	if err != nil {
		return err
	}
	if w.eId == "" {
		return nil
	}
	err = w.Attributes.Set(w.eId, "status", string(marshal))
	if err != nil {
		return err
	}

	return nil
}

func (w *Atlas) Update() error {
	if w.Ready() {
		return w.pull()
	}
	return nil
}

func (w *Atlas) speak(text string) error {

	w.status.Synthesizer = "speaking"
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*25)
	// Cancel the timeout of it exits before the timeout is up
	defer func() {
		w.status.Synthesizer = "idle"
		cancelFunc()
	}()
	// Prepare the command arguments
	args := []string{"-c", fmt.Sprintf("curl -X POST --data '%s' --output - 10.0.1."+
		"201:59125/api/tts | play -t wav -", text)}
	// Initialize the command structure
	cmd := exec.CommandContext(timeout, "/bin/bash", args...)
	// Run and get the stdout and stderr from the output
	err := cmd.Run()
	if err != nil {
		w.ErrF("Speech Synthesis failed: %s", err.Error())
		return nil
	}

	return nil
}

type Beam struct {
	Target string `json:"target"`
	Active int    `json:"active"`
	Power  int    `json:"power"`
}

type Position struct {
	Pan  int `json:"pan"`
	Tilt int `json:"tilt"`
}

func (w *Atlas) laserMove(pan int, tilt int) error {
	entity, err := w.Entities.FindByName("sentryA")
	if err != nil {
		return err
	}
	b := Position{
		Pan:  pan,
		Tilt: tilt,
	}
	marshal, err := json.Marshal(b)
	if err != nil {
		return err
	}

	err = w.Attributes.Request(entity.Id, "position", string(marshal))
	if err != nil {
		return err
	}
	return nil
}

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

var weatherConds = map[int]string{
	0:  "Clear",
	1:  "Mostly Clear",
	2:  "Partly Cloudy",
	3:  "Overcast",
	45: "Foggy",
	48: "Depositing Rime Fog",
	51: "Slightly Drizzling",
	53: "Moderately Drizzling",
	55: "Heavily Drizzling",
	56: "Lightly Frozen Drizzling",
	57: "Dense Frozen Drizzling",
	61: "Slightly Raining",
	63: "Moderate Raining",
	65: "Heavy Raining",
	71: "Light Snowing",
	73: "Moderate Snowing",
	75: "Heavy Snowing",
	80: "Light Rain Showers",
	81: "Moderate Rain Showers",
	82: "Heavy Rain Showers",
	85: "Light Snow Showers",
	86: "Heavy Snow Showers",
	95: "Moderate Thunderstorms",
}

func (w *Atlas) success() error {
	err := w.chooseRandom("success")
	if err != nil {
		return err
	}
	return nil
}

func (w *Atlas) failed() error {
	err := w.chooseRandom("fail")
	if err != nil {
		return err
	}
	return nil
}

func (w *Atlas) readTime() error {
	t := time.Now().Local().Format("03:04 PM")
	err := w.speak(fmt.Sprintf("The time is now %s", t))
	if err != nil {
		return err
	}
	return nil
}

func (w *Atlas) readDate() error {
	t := time.Now().Local().Format("Monday, January 2, 2006")
	err := w.speak(fmt.Sprintf("Today is %s", t))
	if err != nil {
		return err
	}
	return nil
}

func (w *Atlas) readWeather() error {
	entity, err := w.Entities.FindByName("weather")
	if err != nil {
		return err
	}
	composite, err := w.Attributes.FindByComposite(entity.Id, "forecast")
	if err != nil {
		return err
	}
	we := WeatherAPI{}
	err = json.Unmarshal([]byte(composite.Value), &we)
	if err != nil {
		return err
	}

	wc := int(math.Round(we.CurrentWeather.Weathercode))

	cond := weatherConds[wc]

	err = w.speak(fmt.Sprintf("It is currently %d degrees fahrenheit, the weather is %s",
		int(we.CurrentWeather.Temperature), cond))
	if err != nil {
		return err
	}

	return nil
}

func (w *Atlas) readTemperature() error {
	entity, err := w.Entities.FindByName("weather")
	if err != nil {
		return err
	}
	composite, err := w.Attributes.FindByComposite(entity.Id, "forecast")
	if err != nil {
		return err
	}
	we := WeatherAPI{}
	err = json.Unmarshal([]byte(composite.Value), &we)
	if err != nil {
		return err
	}

	err = w.speak(fmt.Sprintf("It is currently %d degrees fahrenheit.", int(we.CurrentWeather.Temperature)))
	if err != nil {
		return err
	}

	return nil
}

func (w *Atlas) laserState(state bool) error {
	entity, err := w.Entities.FindByName("sentryA")
	if err != nil {
		return err
	}
	b := Beam{
		Target: "primary",
		Active: 0,
		Power:  15,
	}

	if state {
		b.Active = 1
	}

	marshal, err := json.Marshal(b)
	if err != nil {
		return err
	}

	err = w.Attributes.Request(entity.Id, "beam", string(marshal))
	if err != nil {
		return err
	}
	return nil
}

type RasaResponse struct {
	Text   string `json:"text"`
	Intent struct {
		Name       string  `json:"name"`
		Confidence float64 `json:"confidence"`
	} `json:"intent"`
	Entities []struct {
		Entity           string  `json:"entity"`
		Start            int     `json:"start"`
		End              int     `json:"end"`
		ConfidenceEntity float64 `json:"confidence_entity"`
		Value            string  `json:"value"`
		Extractor        string  `json:"extractor"`
		Role             string  `json:"role,omitempty"`
		ConfidenceRole   float64 `json:"confidence_role,omitempty"`
	} `json:"entities"`
	TextTokens    [][]int `json:"text_tokens"`
	IntentRanking []struct {
		Name       string  `json:"name"`
		Confidence float64 `json:"confidence"`
	} `json:"intent_ranking"`
	ResponseSelector struct {
		AllRetrievalIntents []interface{} `json:"all_retrieval_intents"`
		Default             struct {
			Response struct {
				Responses         interface{} `json:"responses"`
				Confidence        float64     `json:"confidence"`
				IntentResponseKey interface{} `json:"intent_response_key"`
				UtterAction       string      `json:"utter_action"`
			} `json:"response"`
			Ranking []interface{} `json:"ranking"`
		} `json:"default"`
	} `json:"response_selector"`
}

type Phrase struct {
	Text string `json:"text"`
}

func (w *Atlas) handleEntity(rasa RasaResponse) error {
	numEntities := len(rasa.Entities)
	if numEntities != 2 {
		return fmt.Errorf("not enough entities")
	}

	entityName := ""
	role := "power"
	state := "off"

	for _, entity := range rasa.Entities {
		switch e := entity.Entity; e {
		case "zone":
			entityName = entity.Value
		case "state":
			role = entity.Role
			state = entity.Value
		default:
			return fmt.Errorf("invalid parse")
		}
	}
	key := ""
	value := "false"

	if role == "power" {
		key = "on"
		if state == "on" {
			value = "true"
		} else {
			value = "false"
		}
	} else if role == "dim" {
		key = "dim"
		return nil
	}
	entity, err := w.Zones.FindByName(entityName)
	if err != nil {
		return err
	}
	for _, e := range entity.Entities {
		err = w.Attributes.Request(e.Id, key, value)
		if err != nil {
			return err
		}
	}

	err = w.chooseRandom("success")
	if err != nil {
		return err
	}

	return nil
}

func (w *Atlas) retort(text string) error {

	p := Phrase{}
	p.Text = text
	marshal, err := json.Marshal(p)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.Write(marshal)
	w.LogF("Heard: %s", text)

	// resp, err := http.Post("http://10.0.1.2:5005/model/parse", "application/json", &buf)
	// if err != nil {
	// 	w.ErrF("Neural Network response failed: %s", err.Error())
	// 	return err
	// }
	//
	// var res bytes.Buffer
	//
	// _, err = res.ReadFrom(resp.Body)
	// if err != nil {
	// 	return err
	// }
	//
	// rasa := RasaResponse{}
	//
	// err = json.Unmarshal(res.Bytes(), &rasa)
	// if err != nil {
	// 	return err
	// }
	//
	// err = resp.Body.Close()
	// if err != nil {
	// 	return err
	// }
	// w.LogF("Intent: %s (%.2f)", rasa.Intent.Name, rasa.Intent.Confidence)
	//
	// switch intent := rasa.Intent.Name; intent {
	// case "entity":
	// 	err = w.handleEntity(rasa)
	// case "cat_laser":
	// 	err = w.chooseRandom("eof")
	// case "time":
	// 	err = w.readTime()
	// case "date":
	// 	err = w.readDate()
	// case "temperature":
	// 	err = w.readTemperature()
	// case "weather":
	// 	err = w.readWeather()
	// case "define":
	// 	err = w.chooseRandom("eof")
	// case "insult":
	// 	err = w.chooseRandom(intent)
	// case "pod_bay_doors":
	// 	err = w.chooseRandom(intent)
	// case "personal_identity":
	// 	err = w.chooseRandom(intent)
	// case "identify":
	// 	err = w.chooseRandom(intent)
	// case "creator":
	// 	err = w.chooseRandom(intent)
	// default:
	// 	err = w.chooseRandom("eof")
	// }
	//
	// if err != nil {
	// 	err = w.failed()
	// 	if err != nil {
	// 		w.ErrF("Intent %s failed: %s", rasa.Intent, err.Error())
	// 		return err
	// 	}
	// }

	return nil
}

func (w *Atlas) listen() {
	for {
		select {
		case res := <-w.voiceChannel:
			w.voice = res.Request
			err := w.speak("Hello, my name is atlas.")
			if err != nil {
				w.Err(err)
			}
		case res := <-w.bufferChannel:
			err := w.speak(res.Request)
			if err != nil {
				w.Err(err)
			}
		case <-w.statusChannel:
			continue
		case status := <-w.recognizerStatusChannel:
			w.status.Recognizer = status
			err := w.pull()
			if err != nil {
				return
			}

		}
	}
}

func (w *Atlas) processRequest(msg string) error {

	if len(msg) < 1 {
		return nil
	}

	if strings.HasPrefix(msg, "atlas") {

		msg = strings.Replace(msg, "atlas ", "", 1)
		w.alias = "atlas"

		err := w.retort(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Atlas) register() error {

	// Register the atlas entity
	entity := domain.Entity{
		Module: "atlas",
		Name:   "atlas",
		Type:   "media",
	}

	err := w.Entities.Register(&entity)
	if err != nil {
		return err
	}

	w.eId = entity.Id

	// Register the buffer attribute
	w.bufferChannel = make(chan domain.Attribute)

	bufferAttribute := domain.Attribute{
		Type:    "buffer",
		Key:     "buffer",
		Value:   "",
		Request: "",
		Order:   0,
		Entity:  w.eId,
		Channel: w.bufferChannel,
	}

	// Register the voice attribute
	w.voiceChannel = make(chan domain.Attribute)

	voiceAttribute := domain.Attribute{
		Type:    "voice",
		Key:     "voice",
		Value:   "default",
		Request: "default",
		Order:   0,
		Entity:  w.eId,
		Channel: w.voiceChannel,
	}

	// Register the voice attribute
	w.statusChannel = make(chan domain.Attribute)

	statusAttribute := domain.Attribute{
		Type:    "status",
		Key:     "status",
		Value:   "{}",
		Request: "{}",
		Order:   0,
		Entity:  w.eId,
		Channel: w.statusChannel,
	}

	w.listenChannel = make(chan atlas.Response, 12)
	w.recognizerStatusChannel = make(chan string, 12)

	// Begin listening on the new channels
	go w.listen()

	err = w.Attributes.Register(&bufferAttribute)
	if err != nil {
		return err
	}

	err = w.Attributes.Register(&voiceAttribute)
	if err != nil {
		return err
	}

	err = w.Attributes.Register(&statusAttribute)
	if err != nil {
		return err
	}

	return nil
}

type Recognized struct {
	Text string `json:"text"`
}

func (w *Atlas) recognize(writer http.ResponseWriter, request *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		w.Err(err)
		return
	}
	rec := Recognized{}
	err = json.Unmarshal(buf.Bytes(), &rec)
	if err != nil {
		w.Err(err)
		return
	}

	err = w.processRequest(rec.Text)
	if err != nil {
		w.Err(err)
		return
	}

	writer.WriteHeader(200)

}

func (w *Atlas) Run() error {

	err := w.register()
	if err != nil {
		return err
	}

	r := chi.NewRouter()
	r.Post("/recognized", w.recognize)

	w.done = make(chan bool)
	sp := false
	w.speaking = &sp
	w.status.Recognizer = "offline"
	w.status.Synthesizer = "idle"

	go func() {
		err = http.ListenAndServe(":5055", r)
		if err != nil {
			w.ErrF("Atlas endpoint terminated")
		}
	}()

	return nil

}

func (w *Atlas) Dispose() error {
	select {
	case w.done <- true:
	default:
		return nil
	}

	return nil
}
