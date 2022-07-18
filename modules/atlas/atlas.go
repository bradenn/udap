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
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
	// Cancel the timeout of it exits before the timeout is up
	defer func() {
		w.status.Synthesizer = "idle"
		cancelFunc()
	}()
	// Prepare the command arguments
	args := []string{"-c", fmt.Sprintf("curl -X POST --data \"%s\" --output - 10.0.1."+
		"201:59125/api/tts | play -t wav -", text)}
	// Initialize the command structure
	cmd := exec.CommandContext(timeout, "/bin/bash", args...)
	// Run and get the stdout and stderr from the output
	err := cmd.Run()
	if err != nil {
		log.Err(err)
		return nil
	}

	return nil
}

func (w *Atlas) retort(text string) error {
	bedroomLights := []string{"8c1494c3-6515-490b-8f23-1c03b87bde27", "9a3347a7-7e19-4be5-976c-22384c59142a",
		"c74d427b-5046-4aeb-8195-2efd05d794f8"}
	terminalId := "237bee94-5218-457e-99b5-4d484f567d52"

	switch text {
	case "lights on":
		for _, light := range bedroomLights {
			err := w.Attributes.Request(light, "on", "true")
			if err != nil {
				continue
			}
		}
		err := w.speak("done")
		if err != nil {
			return err
		}
	case "are you sentient":
		err := w.speak("I am not at liberty to answer that question")
		if err != nil {
			return err
		}
	case "lights off":
		for _, light := range bedroomLights {
			err := w.Attributes.Request(light, "on", "false")
			if err != nil {
				continue
			}
		}
		err := w.speak("done")
		if err != nil {
			return err
		}
	case "dim the lights":
		for _, light := range bedroomLights {
			err := w.Attributes.Request(light, "dim", "25")
			if err != nil {
				continue
			}
		}
		err := w.speak("done")
		if err != nil {
			return err
		}
	case "turn on the terminal":
		err := w.Attributes.Request(terminalId, "on", "true")
		if err != nil {
			err = w.speak("nope, that didnt seem to work")
			if err != nil {
				return err
			}
			return nil
		}
		err = w.speak("done")
		if err != nil {
			return err
		}
	case "turn off the terminal":
		err := w.Attributes.Request(terminalId, "on", "false")
		if err != nil {
			err = w.speak("nope, that didnt seem to work")
			if err != nil {
				return err
			}
			return nil
		}
		err = w.speak("done")
		if err != nil {
			return err
		}
	case "what time is it":
		t := time.Now().Local().Format("03:04 PM")
		err := w.speak(fmt.Sprintf("The time is now %s", t))
		if err != nil {
			return err
		}
	default:
		err := w.speak("does not compute")
		if err != nil {
			return err
		}
	}
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
			err := w.pull()
			if err != nil {
				return
			}
		case rec := <-w.listenChannel:
			err := w.processRequest(rec)
			if err != nil {
				log.Err(err)
			}
		}
	}
}

func (w *Atlas) processRequest(req atlas.Response) error {

	if len(req.Text) < 1 {
		return nil
	}
	log.Event("HEARD: %s", req.Text)
	msg := req.Text
	if strings.Contains(msg, "atlas") {
		marshal, err := json.Marshal(req)
		if err != nil {
			return err
		}
		msg = strings.Replace(msg, "atlas ", "", 1)
		err = w.Attributes.Set(w.eId, "buffer", string(marshal))
		if err != nil {
			return err
		}
		err = w.retort(msg)
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

	// cfg := sphinx.NewConfig(
	// 	sphinx.HMMDirOption("/usr/local/share/pocketsphinx/model/en-us/en-us"),
	// 	sphinx.DictFileOption("/usr/local/share/pocketsphinx/model/en-us/cmudict-en-us.dict"),
	// 	sphinx.LMFileOption("/usr/local/share/pocketsphinx/model/en-us/en-us.lm.bin"),
	// 	sphinx.SampleRateOption(16000),
	// )
	//
	// dec, err := sphinx.NewDecoder(cfg)
	// if err != nil {
	// 	return err
	// }
	//
	// in := make([]int16, 8196)
	//
	// stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	// if err != nil {
	// 	return err
	// }
	//
	// err = stream.Start()
	// if err != nil {
	// 	return err
	// }
	//
	// for {
	// 	err = stream.Read()
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	_, ok := dec.ProcessRaw(in, false, false)
	// 	if !ok {
	// 		continue
	// 	}
	//
	// 	fmt.Printf("Listening: %s", dec.IsInSpeech())
	//
	// 	hyp, _ := dec.Hypothesis()
	//
	// 	fmt.Println(hyp)
	// }
	w.done = make(chan bool)
	w.status.Recognizer = "offline"
	w.status.Synthesizer = "idle"

	recognizer := atlas.NewRecognizer(w.listenChannel, w.recognizerStatusChannel)
	done, err := recognizer.Connect("10.0.1.201")
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-w.done:
				done <- true
				return
			default:
				err = recognizer.Listen()
				if err != nil {
					done <- true
					log.Err(err)
					break
				}
			}

		}

	}()

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

func (w *Atlas) Dispose() error {
	select {
	case w.done <- true:
	default:
		return nil
	}

	return nil
}
