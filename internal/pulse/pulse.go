// Copyright (c) 2022 Braden Nicholson

package pulse

import (
	"runtime"
	"strings"
	"sync"
	"time"
	"udap/internal/log"
)

var Timings Timing

func init() {
	Timings = Timing{}
	Timings.history = map[uintptr]Proc{}
	Timings.waiting = map[uintptr]Proc{}
	Timings.handler = make(chan Proc, 1)
	go Timings.handle()
}

type Proc struct {
	Pointer   uintptr   `json:"pointer"`
	Name      string    `json:"name"`
	Start     time.Time `json:"start"`
	Stop      time.Time `json:"stop"`
	Delta     int       `json:"delta"`
	Frequency int       `json:"frequency"`
	Complete  bool      `json:"complete"`
	Depth     int       `json:"depth"`
}

type Timing struct {
	waiting map[uintptr]Proc
	handler chan Proc
	mt      sync.Mutex
	history map[uintptr]Proc
}

func (h *Timing) Timings() (a map[uintptr]Proc) {
	a = map[uintptr]Proc{}
	h.mt.Lock()
	for i, u := range h.history {
		a[i] = u
	}
	h.mt.Unlock()
	return
}

func (h *Timing) handle() {
	h.mt = sync.Mutex{}
	for proc := range h.handler {
		if !proc.Complete {
			h.waiting[proc.Pointer] = proc
			continue
		}

		resident := h.waiting[proc.Pointer]
		resident.Stop = time.Now()
		resident.Complete = true
		resident.Delta = int(time.Since(resident.Start).Nanoseconds())
		h.mt.Lock()
		h.history[resident.Pointer] = resident
		h.mt.Unlock()
		delete(h.waiting, resident.Pointer)
	}
}

func (h *Timing) beginFixed(freq int, rf *runtime.Func) {
	proc := Proc{}
	proc.Frequency = freq
	proc.Pointer = rf.Entry()
	proc.Start = time.Now()
	proc.Complete = false
	split := strings.Split(rf.Name(), ".")
	proc.Name = strings.Join(split[1:], " ")
	h.handler <- proc
}

func (h *Timing) begin(rf *runtime.Func) {
	proc := Proc{}
	proc.Frequency = 2000
	proc.Pointer = rf.Entry()
	proc.Start = time.Now()
	proc.Complete = false
	split := strings.Split(rf.Name(), ".")
	proc.Name = strings.Join(split[1:], " ")
	h.handler <- proc
}

func (h *Timing) end(rf *runtime.Func) error {
	proc := Proc{}
	proc.Pointer = rf.Entry()
	proc.Complete = true
	h.handler <- proc
	return nil
}

func Fixed(ms int) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		Timings.beginFixed(ms, details)
	}
}

func Begin() {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		Timings.begin(details)
	}
}

func End() {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		err := Timings.end(details)
		if err != nil {
			log.Err(err)
		}
	}
}
