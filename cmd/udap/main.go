// Copyright (c) 2021 Braden Nicholson

package main

import (
	"os"
	"udap/internal/log"
	"udap/internal/udap"
)

func main() {

	_, err := udap.Run()
	if err != nil {
		log.ErrF(err, "udap exited due to an error:")
	}
	// Exit normally
	os.Exit(0)
}
