// Copyright (c) 2021 Braden Nicholson

package main

import (
	"udap/internal/log"
	"udap/internal/udap"
)

func main() {
	err := udap.Run()
	if err != nil {
		log.ErrF(err, "UDAP exited due to an unknown error:")
	}
}
