// Copyright (c) 2021 Braden Nicholson

package main

import (
	"os"
	"udap/internal/server"
	"udap/internal/store"
	"udap/internal/udap"
)

func main() {

	// Attempt to recover the program during a panic
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.ErrF(fmt.Errorf("%s\n", r), "panic! Shutting down gracefully.")
	// 		// Initialize all modules
	// 		// err := udap.Cleanup()
	// 		// if err != nil {
	// 		// 	panic(err)
	// 		// }
	// 		log.Log("Good bye.")
	// 	}
	// }()

	u := udap.New()
	u.Add(&store.Database{}, &udap.Build{}, &server.Server{}, &server.Runtime{})
	err := u.Go()
	if err != nil {
		return
	}

	// // Initialize all modules
	// err := udap.Setup()
	// if err != nil {
	// 	panic(err)
	// }
	// // Run udap
	// err = udap.Run()
	// if err != nil {
	// 	return
	// }
	// Exit normally
	os.Exit(0)
}
