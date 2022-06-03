// Copyright (c) 2021 Braden Nicholson

package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kevwan/chatbot/bot"
	"net/url"
	"os/exec"
	"strings"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/pkg/plugin"
)

var Module Atlas

type Atlas struct {
	plugin.Module
	eId        string
	lastSpoken string
	chatbot    *bot.ChatBot
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

func init() {
	config := plugin.Config{
		Name:        "atlas",
		Type:        "module",
		Description: "General AI",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}
	Module.Config = config
}

func (w *Atlas) Setup() (plugin.Config, error) {
	err := w.UpdateInterval(2000)
	if err != nil {
		return plugin.Config{}, err
	}
	return w.Config, nil
}

func (w *Atlas) pull() error {

	return nil
}

func (w *Atlas) Update() error {
	if time.Since(w.Module.LastUpdate) >= time.Second*2 {
		w.Module.LastUpdate = time.Now()
		return w.pull()
	}
	return nil
}

func (w *Atlas) speak(text string) error {

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
	// Cancel the timeout of it exits before the timeout is up
	defer cancelFunc()
	// Prepare the command arguments
	args := []string{"-t", text, "-voice", "./pkg/mimic/mycroft_voice_4.0.flitevox"}
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

func (w *Atlas) register() error {
	entity := domain.Entity{
		Module: "atlas",
		Name:   "atlas",
		Type:   "media",
	}

	err := w.Entities.Register(&entity)
	if err != nil {
		return err
	}

	listenBuffer := domain.Attribute{
		Type:    "buffer",
		Key:     "buffer",
		Value:   "",
		Request: "",
		Order:   0,
		Entity:  entity.Id,
		Channel: make(chan domain.Attribute),
	}

	w.eId = entity.Id

	go func() {
		for attribute := range listenBuffer.Channel {
			fmt.Println("Atlas hears: " + attribute.Value)
		}
	}()

	err = w.Attributes.Register(&listenBuffer)
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

	u := url.URL{Scheme: "ws", Host: "localhost" + ":" + "2700", Path: ""}

	// Opening websocket connection
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer c.Close()

	for {
		msg := Message{}
		err = c.ReadJSON(&msg)
		if err != nil {
			return err
		}
		err = w.Attributes.Set(w.eId, "buffer", msg.Text)
		if err != nil {
			return err
		}

		if strings.Contains(msg.Text, "atlas") {
			if strings.HasPrefix(msg.Text, "the") {
				msg.Text = strings.Replace(msg.Text, "the ", "", 1)
			}
			msg.Text = strings.Replace(msg.Text, "atlas ", "", 1)
			err = w.retort(msg.Text)
			if err != nil {
				return err
			}
		}

		log.Event("ATLAS: %s", msg.Text)
	}

}
