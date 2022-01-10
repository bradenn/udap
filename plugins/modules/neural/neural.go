// Copyright (c) 2021 Braden Nicholson

package main

import (
	"encoding/json"
	"github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
	"time"
	"udap/internal/log"
	"udap/internal/models"
	"udap/pkg/plugin"
)

var Module Neural

type Neural struct {
	plugin.Module
	lastRun time.Time
}

func init() {
	config := plugin.Config{
		Name:        "neural",
		Type:        "module",
		Description: "Neural Network Integration",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}
	Module.Config = config
}

// hour, day, week, current state

func (n *Neural) Suggest(id string) error {
	var examples []training.Example
	var logs []models.Log

	a, err := n.Bond.SendId("entity", "logs", id, nil)
	if err != nil {
		return err
	}

	logs = a.([]models.Log)

	for _, lg := range logs {
		day := float64(lg.CreatedAt.Hour() / 23.0) // 0 - 23 -> 0 -> 1
		level := float64(lg.Level / 100)
		cct := float64((lg.CCT - 2000) / 6000)
		power := float64(0)
		if lg.Power == "on" && lg.Level > 0 {
			power = 1
		}
		output := []float64{power}
		input := []float64{day, level, cct, power}

		example := training.Example{Input: input, Response: output}
		examples = append(examples, example)
	}

	data := training.Examples(examples)

	ne := deep.NewNeural(&deep.Config{
		/* Input dimensionality */
		Inputs: 4,
		/* Two hidden layers consisting of two neurons each, and a single output */
		Layout: []int{4, 4, 1},
		/* Activation functions: Sigmoid, Tanh, ReLU, Linear */
		Activation: deep.ActivationSigmoid,
		/* Determines output layer activation & loss function:
		ModeRegression: linear outputs with MSE loss
		ModeMultiClass: softmax output with Cross Entropy loss
		ModeMultiLabel: sigmoid output with Cross Entropy loss
		ModeBinary: sigmoid output with binary CE loss */
		Mode: deep.ModeBinary,
		/* Weight initializers: {deep.NewNormal(μ, σ), deep.NewUniform(μ, σ)} */
		Weight: deep.NewNormal(1.0, 0.0),
		/* Apply bias */
		Bias: true,
	})

	// params: learning rate, momentum, alpha decay, nesterov
	optimizer := training.NewSGD(0.05, 0.1, 1e-6, true)
	// params: optimizer, verbosity (print stats at every 50th iteration)
	trainer := training.NewTrainer(optimizer, 2000)

	tra, heldout := data.Split(0.5)
	trainer.Train(ne, tra, heldout, 2000) // training, validation, iterations

	a, err = n.Bond.SendId("entity", "pull", id, nil)
	if err != nil {
		return err
	}

	sta := a.(json.RawMessage)

	state := models.LightState{}
	state.Parse(sta)

	day := float64(time.Now().Hour()) / 23.0 // 0 - 23 -> 0 -> 1
	level := float64(state.Level) / 100.0
	cct := float64(state.CCT-2000) / 6000.0

	power := 0.0
	if state.Power == "on" {
		power = 1.0
	}

	input := []float64{day, level, cct, power}
	pd := ne.Predict(input)
	if pd[0] >= 0.5 {
		state.Power = "on"
	} else {
		state.Power = "off"
	}

	_, err = n.Bond.SendId("entity", "predict", id, state.JSON())
	if err != nil {
		return err
	}
	return nil
}

func (n *Neural) Setup() (plugin.Config, error) {
	n.lastRun = time.Now()
	return n.Config, nil
}

func (n *Neural) Update() error {
	ents, err := n.Bond.SendId("entity", "compile", "", nil)
	if err != nil {
		return err
	}

	entities := ents.([]models.Entity)
	for _, entity := range entities {
		go func(id string) {
			err = n.Suggest(id)
			if err != nil {
				return
			}
		}(entity.Id)
	}
	return nil
}

func (n *Neural) Run() error {
	for {
		err := n.Update()
		if err != nil {
			log.Err(err)
		}
		time.Sleep(time.Minute * 10)
	}

}
