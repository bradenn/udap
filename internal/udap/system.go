// Copyright (c) 2021 Braden Nicholson

package udap

import (
	"runtime"
	"udap/internal/log"
)

type Environment struct {
}

type System struct {
	Cores   int `json:"cores"`
	Threads int `json:"threads"`
}
type Memory struct {
}

func GetMem() {
	s := System{}
	s.Threads = runtime.NumGoroutine()
	s.Cores = runtime.NumCPU()
	var sr []runtime.StackRecord
	runtime.Gosched()

	log.Log("%s", sr)
}
