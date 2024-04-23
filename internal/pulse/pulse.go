// Copyright (c) 2022 Braden Nicholson

package pulse

import (
	"sync"
	"time"
)

var Timings Timing

func init() {
	Timings = Timing{}
	Timings.mt = sync.RWMutex{}
	Timings.historyMutex = sync.RWMutex{}
	Timings.history = map[string]Proc{}
	Timings.waiting = map[string]Proc{}
	Timings.handler = make(chan Proc, 32)
	go Timings.handle()
}

type Timing struct {
	waiting      map[string]Proc
	handler      chan Proc
	mt           sync.RWMutex
	historyMutex sync.RWMutex
	history      map[string]Proc
}

type Proc struct {
	Pointer   string    `json:"pointer"`
	Name      string    `json:"name"`
	Start     time.Time `json:"start"`
	StartNano int64     `json:"startNano"`
	Stop      time.Time `json:"stop"`
	StopNano  int64     `json:"stopNano"`
	Delta     int       `json:"delta"`
	Frequency int64     `json:"frequency"`
	Complete  bool      `json:"complete"`
	Depth     int       `json:"depth"`
}

// Timings returns the present timings manifest
func (h *Timing) Timings() (a map[string]Proc) {
	a = map[string]Proc{}
	h.historyMutex.Lock()
	for i, u := range h.history {

		if time.Since(u.Stop).Seconds() < 60 {
			a[i] = u
		} else {
			//_, ok := h.history[i]
			//if ok {
			//	delete(h.history, i)
			//}
		}
	}
	h.historyMutex.Unlock()
	//fmt.Println("Sent timings")
	return
}

func (h *Timing) AllTimings() (a map[string]Proc) {
	a = map[string]Proc{}
	h.historyMutex.Lock()
	for i, u := range h.history {
		a[i] = u
	}
	h.historyMutex.Unlock()
	//fmt.Println("Sent timings")
	return
}

func (h *Timing) await(proc Proc) {
	h.mt.Lock()
	h.waiting[proc.Pointer] = proc
	h.mt.Unlock()
}

// accept processes a proc request, the proc will either be added to the waiting queue or be resolved
func (h *Timing) accept(proc Proc) {
	// Send to waiting queue
	if !proc.Complete {
		h.await(proc)
		return
	}
	h.mt.RLock()
	// Resolve waiting process
	resident := h.waiting[proc.Pointer]
	h.mt.RUnlock()

	// Mark the process as complete
	resident.Complete = true
	resident.Stop = proc.Stop
	resident.StopNano = proc.StopNano
	// Compute the delta in nanoseconds
	resident.Delta = int(time.Since(resident.Start).Nanoseconds())
	// Lock the mutex and record the record
	h.historyMutex.Lock()
	current, ok := h.history[resident.Pointer]
	if ok {
		resident.Frequency = resident.StartNano - current.StopNano
	}
	h.history[resident.Pointer] = resident
	h.historyMutex.Unlock()
	// Delete the resident from the waiting queue
	h.mt.Lock()
	delete(h.waiting, resident.Pointer)
	h.mt.Unlock()
}

// handle accepts incoming timing requests
func (h *Timing) handle() {
	for {
		select {
		// Accept any incoming requests, handle them in a new thread
		case proc := <-h.handler:
			h.accept(proc)
		}
	}
}

// begin starts a process timing request
func (h *Timing) begin(name string, pointer string) {
	// Record the start time
	start := time.Now()
	// Create the process request
	proc := Proc{
		Name:      name,
		Pointer:   pointer,
		Start:     start,
		StartNano: start.UnixNano(),
		Complete:  false,
	}
	// Send the request to the handler mux
	h.handler <- proc
}

// end resolves a process timing request
func (h *Timing) end(name string, pointer string) {
	// Record the stop time
	stop := time.Now()
	// Create the process request
	proc := Proc{
		Name:     name,
		Pointer:  pointer,
		Complete: true,
		Stop:     stop,
		StopNano: stop.UnixNano(),
	}
	// Send the request to the handler mux
	h.handler <- proc

}

func Begin(ref string) {
	Timings.begin(ref, ref)
}

func End(ref string) {
	Timings.end(ref, ref)
}
