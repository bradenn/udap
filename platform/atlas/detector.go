// Copyright (c) 2022 Braden Nicholson

package atlas

import (
	"github.com/mjibson/go-dsp/fft"
	"math"
)

type Detector struct {
	samples  []complex128
	fft      []complex128
	spectrum []float64
	previous []float64
}

// Mostly copied from https://github.com/airikbeeso/hear/blob/master/vad.go

func NewDetector(width int) *Detector {
	return &Detector{
		samples:  make([]complex128, width),
		spectrum: make([]float64, width/2+1),
		previous: make([]float64, width/2+1),
	}
}

func (d *Detector) Detect(samples []int16) float64 {
	for i, sample := range samples {
		d.samples[i] = complex(float64(sample), 0)
	}
	d.fft = fft.FFT(d.samples)
	copy(d.spectrum, d.previous)

	for i := range d.spectrum {
		local := d.fft[i]
		d.spectrum[i] = math.Sqrt(real(local)*real(local) + imag(local)*imag(local))
	}

	var delta float64

	for i, sample := range d.spectrum {
		delta += sample - d.previous[i]
	}

	return delta
}

func (d *Detector) FFT() []complex128 {
	return d.fft
}
