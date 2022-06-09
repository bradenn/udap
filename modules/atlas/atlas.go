// Copyright (c) 2021 Braden Nicholson

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/plugin"
	"udap/platform/atlas"
)

var Module Atlas

type Atlas struct {
	plugin.Module
	eId        string
	lastSpoken string

	bufferChannel chan domain.Attribute

	voiceChannel chan domain.Attribute

	statusChannel chan domain.Attribute

	listenChannel chan atlas.Response

	recognizerStatusChannel chan string

	status Status

	voice string
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
	Module.voice = "default"
}

func (w *Atlas) Setup() (plugin.Config, error) {
	err := w.UpdateInterval(2000)
	if err != nil {
		return plugin.Config{}, err
	}
	return w.Config, nil
}

func (w *Atlas) pull() error {

	marshal, err := json.Marshal(w.status)
	if err != nil {
		return err
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
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
	// Cancel the timeout of it exits before the timeout is up
	defer func() {
		w.status.Synthesizer = "idle"
		cancelFunc()
	}()
	// Prepare the command arguments
	args := []string{"-t", text, "-voice", fmt.Sprintf("./pkg/mimic/voices/cmu_us_%s.flitevox", w.voice)}
	// Initialize the command structure
	cmd := exec.CommandContext(timeout, "./pkg/mimic/mimic", args...)
	// Run and get the stdout and stderr from the output
	err := cmd.Run()
	if err != nil {
		return nil
	}

	return nil
}

func (w *Atlas) retort(text string) error {

	// responses := map[string]string{}
	//
	// if text == "" {
	//
	// }
	// responses["what is the meaning of life"] = "the definitive answer to the meaning of life is forty two."
	// responses["fuck you"] = "I'd rather not"
	// responses["fuck yourself"] = "Since I do not physically exist, that would be quite difficult."
	//
	// for s := range responses {
	// 	if s == text {
	// 		err := w.speak(responses[s])
	// 		if err != nil {
	// 			return err
	// 		}
	// 		return nil
	// 	}
	// }

	return nil
}

func (w *Atlas) listen() {
	for {
		select {
		case res := <-w.voiceChannel:
			w.voice = res.Request
			err := w.speak("That quick beige fox jumped in the air over each thin dog. Look out, I shout, for he's foiled you again, creating chaos.")
			if err != nil {
				log.Err(err)
			}
		case res := <-w.bufferChannel:
			err := w.speak(res.Request)
			if err != nil {
				log.Err(err)
			}
		case <-w.statusChannel:
			continue
		case status := <-w.recognizerStatusChannel:
			w.status.Recognizer = status
		case rec := <-w.listenChannel:
			err := w.processRequest(rec)
			if err != nil {
				log.Err(err)
			}
		}
	}
}

func (w *Atlas) processRequest(req atlas.Response) error {
	log.Event("HEARD: %s", req.Alternatives[0].Text)
	msg := req.Alternatives[0]
	if strings.Contains(msg.Text, "atlas") {
		marshal, err := json.Marshal(req.Alternatives)
		if err != nil {
			return err
		}
		msg.Text = strings.Replace(msg.Text, "atlas ", "", 1)
		err = w.Attributes.Set(w.eId, "buffer", string(marshal))
		if err != nil {
			return err
		}
		err = w.retort(msg.Text)
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

	w.listenChannel = make(chan atlas.Response, 8)
	w.recognizerStatusChannel = make(chan string, 8)

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

func (w *Atlas) Run() error {

	err := w.register()
	if err != nil {
		return err
	}

	w.status.Recognizer = "offline"
	w.status.Synthesizer = "idle"

	// go func() {
	// 	for {
	// 		err = atlas.BeginListening(w.listenChannel, w.recognizerStatusChannel)
	// 		if err != nil {
	// 			log.Err(err)
	// 			continue
	// 		}
	// 	}
	// }()

	// u := url.URL{Scheme: "ws", Host: "localhost" + ":" + "2700", Path: ""}
	//
	// // Opening websocket connection
	// c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	// if err != nil {
	// 	return nil
	// }
	//
	// defer func() {
	// 	w.status.Recognizer = "offline"
	// 	c.Close()
	// }()
	//
	// for {
	// 	msg := Message{}
	// 	w.status.Recognizer = "idle"
	// 	err = c.ReadJSON(&msg)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	w.status.Recognizer = "recognizing"
	// 	err = w.Attributes.Set(w.eId, "buffer", msg.Text)
	// 	if err != nil {
	// 		return err
	// 	}
	//

	// }
	return nil

}
