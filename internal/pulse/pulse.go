// Copyright (c) 2022 Braden Nicholson

package pulse

import (
	"sync"
	"time"
	"udap/internal/log"
)

var Timings Timing

func init() {
	Timings = Timing{}
	Timings.mt = sync.RWMutex{}
	Timings.history = map[string]Proc{}
	Timings.waiting = map[string]Proc{}
	Timings.handler = make(chan Proc, 16)
	go Timings.handle()
}

type Timing struct {
	waiting map[string]Proc
	handler chan Proc
	mt      sync.RWMutex
	history map[string]Proc
}

type Proc struct {
	Pointer   string    `json:"pointer"`
	Name      string    `json:"name"`
	Start     time.Time `json:"start"`
	StartNano int64     `json:"startNano"`
	Stop      time.Time `json:"stop"`
	StopNano  int64     `json:"stopNano"`
	Delta     int       `json:"delta"`
	Frequency int       `json:"frequency"`
	Complete  bool      `json:"complete"`
	Depth     int       `json:"depth"`
}

func (h *Timing) Timings() (a map[string]Proc) {
	a = map[string]Proc{}
	h.mt.RLock()
	for i, u := range h.history {
		a[i] = u
	}
	h.mt.RUnlock()
	return
}

func (h *Timing) handle() {

	for proc := range h.handler {
		if !proc.Complete {
			h.waiting[proc.Pointer] = proc
			continue
		}

		resident := h.waiting[proc.Pointer]
		resident.Stop = time.Now()
		resident.StopNano = resident.Stop.UnixNano()
		resident.Complete = true
		resident.Delta = int(time.Since(resident.Start).Nanoseconds())
		h.mt.Lock()
		h.history[resident.Pointer] = resident
		h.mt.Unlock()
		delete(h.waiting, resident.Pointer)
	}
}

func (h *Timing) begin(ref string) {
	proc := Proc{}
	proc.Pointer = ref
	proc.Start = time.Now()
	proc.StartNano = proc.Start.UnixNano()
	proc.Complete = false
	proc.Name = ref
	h.handler <- proc
}

func (h *Timing) end(ref string) error {
	proc := Proc{}
	proc.Pointer = ref
	proc.Complete = true
	h.handler <- proc
	return nil
}

func Begin(ref string) {
	// pc, _, _, ok := runtime.Caller(1)
	// details := runtime.FuncForPC(pc)
	Timings.begin(ref)
	// if ok && details != nil {
	// }
}

func End(ref string) {
	// pc, _, _, ok := runtime.Caller(1)
	// details := runtime.FuncForPC(pc)
	err := Timings.end(ref)
	if err != nil {
		log.Err(err)
	}
	// if ok && details != nil {
	//
	// }
}