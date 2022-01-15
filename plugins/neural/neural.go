// Copyright (c) 2021 Braden Nicholson

package main

import (
	"encoding/json"
	"github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
	"math"
	"sync"
	"time"
	"udap/internal/log"
	"udap/internal/models"
	"udap/pkg/plugin"
)

var Module Neural

type Neural struct {
	plugin.Module
	lastRun time.Time
	trained map[string]*deep.Neural
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
	ne := n.trained[id]
	if ne == nil {
		err := n.Train(id)
		if err != nil {
			return err
		}
		ne = n.trained[id]
	}
	ent, err := n.Entities.Pull(id)
	state := models.LightState{}
	state.Parse(json.RawMessage(ent))

	day := float64(time.Now().Hour()) / 23.0 // 0 - 23 -> 0 -> 1
	cct := float64(state.CCT-2000) / 6000.0
	// level := float64(state.Level) / 100.0

	input := []float64{day, cct}
	pd := ne.Predict(input)
	state.Level = int(math.Round((pd[0]*100.0)/5.0)) * 5
	state.Mode = "brightness"

	log.Log("suggesting")
	_, err = n.Entities.Suggest(id, string(state.JSON()))
	if err != nil {
		return err
	}

	return nil
}

func (n *Neural) Train(id string) error {

	var examples []training.Example
	var logs []models.Log

	a, err := n.Bond.SendId("entity", "logs", id, nil)
	if err != nil {
		return err
	}

	logs = a.([]models.Log)

	for _, lg := range logs {
		day := float64(lg.CreatedAt.Hour()) / 23.0 // 0 - 23 -> 0 -> 1
		level := float64(lg.Level) / 100.0
		cct := float64(lg.CCT-2000) / 6000.0

		output := []float64{level}
		input := []float64{day, cct}

		example := training.Example{Input: input, Response: output}
		examples = append(examples, example)
	}

	data := training.Examples(examples)

	ne := deep.NewNeural(&deep.Config{
		/* Input dimensionality */
		Inputs: 2,
		/* Two hidden layers consisting of two neurons each, and a single output */
		Layout: []int{2, 2, 1},
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
	trainer := training.NewTrainer(optimizer, 5000)

	tra, heldout := data.Split(0.5)
	trainer.Train(ne, tra, heldout, 8000) // training, validation, iterations
	n.trained[id] = ne

	// err = ent.Push(models.State(state.JSON()))
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (n *Neural) Setup() (plugin.Config, error) {
	n.lastRun = time.Now()
	n.trained = map[string]*deep.Neural{}
	return n.Config, nil
}

func (n *Neural) Update() error {
	entities, err := n.Entities.Compile()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}

	for _, entity := range entities {
		if entity.Neural == "suggest" || entity.Neural == "control" {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				err = n.Suggest(id)
				if err != nil {
					return
				}
			}(entity.Id)
		}
	}
	wg.Wait()
	return nil
}

func (n *Neural) Run() error {
	time.Sleep(time.Second * 5)
	// for {
	// 	// err := n.Update()
	// 	// if err != nil {
	// 	// 	log.Err(err)
	// 	// }
	// 	time.Sleep(time.Second * 10)
	// }
	return nil
}
