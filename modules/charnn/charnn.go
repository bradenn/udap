// Copyright (c) 2023 Braden Nicholson

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
	"os"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/plugin"
)

var Module CharNN

// Structure partially inspired by https://sausheong.github.io/posts/how-to-build-a-simple-artificial-neural-network-with-go/

type Network struct {
	inputs        int
	hidden        int
	outputs       int
	hiddenWeights *mat.Dense
	outputWeights *mat.Dense
	learningRate  float64
}

func (n *Network) Inputs() int {
	return n.inputs
}
func (n *Network) Outputs() int {
	return n.outputs
}

func Save(n Network) error {
	hidden, err := os.OpenFile("./modules/charnn/charnn-h.model", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer hidden.Close()
	if err != nil {
		return err
	}

	output, err := os.OpenFile("./modules/charnn/charnn-o.model", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer output.Close()
	if err != nil {
		return err
	}

	_, err = n.hiddenWeights.MarshalBinaryTo(hidden)
	if err != nil {
		return err
	}

	_, err = n.outputWeights.MarshalBinaryTo(output)
	if err != nil {
		return err
	}

	return nil
}

func Load(n *Network) error {

	hidden, err := os.Open("./modules/charnn/charnn-h.model")
	defer hidden.Close()
	if err != nil {
		return err
	}

	output, err := os.Open("./modules/charnn/charnn-o.model")
	defer output.Close()
	if err != nil {
		return err
	}

	n.hiddenWeights.Reset()
	_, err = n.hiddenWeights.UnmarshalBinaryFrom(hidden)
	if err != nil {
		return err
	}

	n.outputWeights.Reset()
	_, err = n.outputWeights.UnmarshalBinaryFrom(output)
	if err != nil {
		return err
	}

	return nil
}

// dot returns a matrix of the dot product of the matrix i and j
func dot(i, j mat.Matrix) mat.Matrix {
	// Get the height of the i matrix
	rows, _ := i.Dims()
	// Get the width of the j matrix
	_, cols := j.Dims()
	// Initialize an output matrix of appropriate size
	output := mat.NewDense(rows, cols, nil)
	// Get the product of the i and j matrix and store it in the output matrix
	output.Product(i, j)
	// Return the new matrix
	return output
}

// apply returns a matrix mutated by a provided function
func apply(fn func(i, j int, v float64) float64, m mat.Matrix) mat.Matrix {
	// Get the size of the provided matrix, m
	rows, cols := m.Dims()
	// Create a new output matrix with the same sizes as m
	output := mat.NewDense(rows, cols, nil)
	// Apply the function to the matrix m and store it in the new matrix
	output.Apply(fn, m)
	// Return the new matrix
	return output
}

// scale multiplies a matrix by a scalar value
func scale(scale float64, m mat.Matrix) mat.Matrix {
	// Get the size of the provided matrix, m
	rows, cols := m.Dims()
	// Create a new output matrix with the same sizes as m
	output := mat.NewDense(rows, cols, nil)
	// Apply the function to the matrix m and store it in the new matrix
	output.Scale(scale, m)
	// Return the new matrix
	return output
}

// add adds two matrices
func add(i, j mat.Matrix) mat.Matrix {
	// Get the size of the provided matrix, m
	rows, cols := i.Dims()
	// Create a new output matrix with the same sizes as m
	output := mat.NewDense(rows, cols, nil)
	// Add the value into the new matrix
	output.Add(i, j)
	// Return the new matrix
	return output
}

// subtract adds two matrices
func subtract(i, j mat.Matrix) mat.Matrix {
	// Get the size of the provided matrix, m
	rows, cols := i.Dims()
	// Create a new output matrix with the same sizes as m
	output := mat.NewDense(rows, cols, nil)
	// Subtract the value into the new matrix
	output.Sub(i, j)
	// Return the new matrix
	return output
}

// addScalar adds a scalar value to a matrix
func addScalar(s float64, m mat.Matrix) mat.Matrix {
	// Get the size of the provided matrix, m
	rows, cols := m.Dims()
	// Create a new array of to represent a matrix
	cache := make([]float64, rows*cols)
	// Fill all values with the array, and in turn, the matrix with the value s to add
	for i := 0; i < rows*cols; i++ {
		// Set each value
		cache[i] = s
	}
	// Create a new matrix with the same sizes as m
	n := mat.NewDense(rows, cols, cache)
	// Return the addition of the scalar matrix and the target
	return add(m, n)
}

// multiply multiplies two matrices
func multiply(i, j mat.Matrix) mat.Matrix {
	// Get the size of the provided matrix, m
	rows, cols := i.Dims()
	// Create a new output matrix with the same sizes as m
	output := mat.NewDense(rows, cols, nil)
	// Multiply the value into the new matrix
	output.MulElem(i, j)
	// Return the new matrix
	return output
}

// generateValues returns an array of uniformly distributed random values
func generateValues(current int, previous int) []float64 {
	dist := distuv.Uniform{
		Min: -1 / math.Sqrt(float64(previous)),
		Max: 1 / math.Sqrt(float64(previous)),
	}
	data := make([]float64, current)
	for i := 0; i < current; i++ {
		data[i] = dist.Rand()
	}
	return data
}

func sigmoid(rows, cols int, value float64) float64 {
	return 1.0 / (1 + math.Exp(-1.0*value))
}

func sigmoidDX(i mat.Matrix) mat.Matrix {
	// Get the height of the i matrix
	rows, _ := i.Dims()
	// Create a matrix array to contain the ones
	output := make([]float64, rows)
	// Fill all values with ones
	for j := range output {
		output[j] = 1
	}
	// Create a matrix from the array
	target := mat.NewDense(rows, 1, output)
	// Calculate the derived sigmoid
	return multiply(i, subtract(target, i)) // i * (1 - i)

}

func NewNetwork(input, hidden, output int, learningRate float64) (network Network) {

	network = Network{
		inputs:       input,
		outputs:      output,
		hidden:       hidden,
		learningRate: learningRate,
	}
	network.hiddenWeights = mat.NewDense(hidden, input, generateValues(input*hidden, input))
	network.outputWeights = mat.NewDense(output, hidden, generateValues(output*hidden, hidden))
	return network
}

func (n *Network) Train(input []float64, output []float64) {

	// Forward Propagation (verbatim from Predict)

	// Create a 1xN sized matrix column to hold the input values
	inputs := mat.NewDense(len(input), 1, input)
	// Multiply the inputs by the hidden layer weights
	hiddenInputs := dot(n.hiddenWeights, inputs)
	// Run the activation function on the values
	hiddenOutputs := apply(sigmoid, hiddenInputs)
	// Multiply the hidden layer weights by the output weights
	finalInputs := dot(n.outputWeights, hiddenOutputs)
	// Apply the activation for the final weights
	finalOutputs := apply(sigmoid, finalInputs)

	// Errors
	// Create a matrix to represent the target output data
	targets := mat.NewDense(len(output), 1, output)
	// Calculate the difference from the prediction and the target result
	outputError := subtract(targets, finalOutputs)
	// Calculate the hidden layer error
	hiddenError := dot(n.outputWeights.T(), outputError)

	// Backpropagation

	// Update the output weights
	n.outputWeights = add(n.outputWeights,
		scale(n.learningRate,
			dot(multiply(outputError, sigmoidDX(finalOutputs)),
				hiddenOutputs.T()))).(*mat.Dense)

	// Update the hidden weights
	n.hiddenWeights = add(n.hiddenWeights,
		scale(n.learningRate,
			dot(multiply(hiddenError, sigmoidDX(hiddenOutputs)),
				inputs.T()))).(*mat.Dense)

}

// Predict runs an array which represent the inputs into the neural network and returns the output values
func (n *Network) Predict(data []float64) mat.Matrix {
	// Create a 1xN sized matrix column to hold the input values
	inputs := mat.NewDense(len(data), 1, data)
	// Multiply the inputs by the hidden layer weights
	hiddenInputs := dot(n.hiddenWeights, inputs)
	// Run the activation function on the values
	hiddenOutputs := apply(sigmoid, hiddenInputs)
	// Multiply the hidden layer weights by the output weights
	finalInputs := dot(n.outputWeights, hiddenOutputs)
	// Apply the activation for the final weights
	finalOutputs := apply(sigmoid, finalInputs)
	// Return the results
	return finalOutputs
}

type CharNN struct {
	plugin.Module
	entityId string
	network  *Network
	receiver chan domain.Attribute
}

func (c *CharNN) mux() {
	for msg := range c.receiver {
		c.handleQuery(msg)
	}
}

type Query struct {
	Stream string `json:"stream"`
}

func (c *CharNN) parseQuery(attribute domain.Attribute) error {
	q := Query{}
	err := json.Unmarshal([]byte(attribute.Request), &q)
	if err != nil {
		return err
	}
	decodeString, err := base64.StdEncoding.DecodeString(q.Stream)
	if err != nil {
		return err
	}

	inputs := make([]float64, c.network.Inputs())
	for i := range inputs {
		x := float64(decodeString[i])
		inputs[i] = (x / 255.0 * 0.99) + 0.01
	}

	output := c.network.Predict(inputs)
	max := 0.0
	value := 0
	for i := 0; i < c.network.Outputs(); i++ {
		a := output.At(i, 0)
		if a > max {
			max = a
			value = i
		}
	}

	err = c.Attributes.Update(c.entityId, "result", fmt.Sprintf("%d", value), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (c *CharNN) handleQuery(msg domain.Attribute) {
	switch msg.Key {
	case "query":
		_ = c.parseQuery(msg)
		break
	default:
		break
	}
}

func (c *CharNN) Setup() (plugin.Config, error) {

	return c.Config, nil
}

func (c *CharNN) Run() error {

	err := c.InitConfig("server", "https://example.com")
	if err != nil {
		log.Err(err)
		return nil
	}
	c.receiver = make(chan domain.Attribute, 8)
	go func() {
		c.mux()
	}()

	entity := domain.Entity{
		Name:   "charnn",
		Type:   "media",
		Module: c.Config.Name,
	}

	err = c.Entities.Register(&entity)
	if err != nil {
		return err
	}

	c.entityId = entity.Id

	query := domain.Attribute{
		Value:     "{}",
		Updated:   time.Now(),
		Request:   "{}",
		Requested: time.Now(),
		Entity:    c.entityId,
		Key:       "query",
		Type:      "media",
		Order:     0,
		Channel:   c.receiver,
	}

	err = c.Attributes.Register(&query)
	if err != nil {
		return err
	}

	results := domain.Attribute{
		Value:     "{}",
		Updated:   time.Now(),
		Request:   "{}",
		Requested: time.Now(),
		Entity:    c.entityId,
		Key:       "result",
		Type:      "media",
		Order:     1,
		Channel:   c.receiver,
	}
	err = c.Attributes.Register(&results)
	if err != nil {
		return err
	}

	net := NewNetwork(784, 250, 10, 0.1)
	c.network = &net

	err = Load(c.network)
	if err != nil {
		log.Err(err)
		return err
	}

	return nil
}

func (c *CharNN) Update() error {
	c.Ready()
	return nil
}

func init() {

	config := plugin.Config{
		Name:        "charnn",
		Type:        "module",
		Description: "Neural Network based drawn character classification",
		Version:     "0.1.0",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}
