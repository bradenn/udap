// Copyright (c) 2021 Braden Nicholson

package main

import (
	"fmt"
	"github.com/Picovoice/porcupine/binding/go/v2"
	pvrecorder "github.com/Picovoice/pvrecorder/sdk/go"
	"log"
	"os"
	"udap/pkg/plugin"
)

var Module Gpt

type Gpt struct {
	plugin.Module
}

func init() {
	config := plugin.Config{
		Name:        "gpt",
		Type:        "module",
		Description: "Gpt heuristics",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func (g *Gpt) Setup() (plugin.Config, error) {

	return g.Config, nil
}

func (g *Gpt) Update() error {
	return nil
}
func (g *Gpt) Run() error {
	return nil
}
func (g *Gpt) run() error {
	p := porcupine.Porcupine{}
	p.AccessKey = os.Getenv("picovoice")
	p.KeywordPaths = []string{"./wakemodel.ppn"}
	err := p.Init()
	if err != nil {
		return err
	}

	defer func(p *porcupine.Porcupine) {
		err = p.Delete()
		if err != nil {

		}
	}(&p)

	recorder := pvrecorder.PvRecorder{
		DeviceIndex:    -1,
		FrameLength:    porcupine.FrameLength,
		BufferSizeMSec: 1000,
		LogOverflow:    0,
	}

	if err := recorder.Init(); err != nil {
		log.Fatalf("Error: %s.\n", err.Error())
	}
	defer recorder.Delete()

	log.Printf("Using device: %s", recorder.GetSelectedDevice())

	if err := recorder.Start(); err != nil {
		log.Fatalf("Error: %s.\n", err.Error())
	}

	for {
		pcm, err := recorder.Read()
		if err != nil {
			log.Fatalf("Error: %s.\n", err.Error())
		}
		keywordIndex, err := p.Process(pcm)
		if err != nil {
			log.Fatal(err)
		}
		if keywordIndex >= 0 {
			fmt.Println("Hello there")
		}
	}

	return nil
}
