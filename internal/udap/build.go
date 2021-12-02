// Copyright (c) 2021 Braden Nicholson

package udap

import (
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"udap/pkg/plugin"
)

type Build struct {
	database *gorm.DB
	wg       *sync.WaitGroup
}

// Dependency is the level at which this service needs to run
func (b *Build) Dependency() (level int) {
	return 1
}

func (b *Build) Name() (name string) {
	return "build"
}

// Load configures and prepares the parent struct for running
func (b *Build) Load() (err error) {
	wg := sync.WaitGroup{}
	plugin.BuildAll("./plugins/modules/", &wg)
	wg.Wait()
	return nil
}

// Run will begin the main-sequence activities of the parent struct
func (b *Build) Run(interface{}) (err error) {

	// No action occurs
	return nil
}

// Cleanup will begin the main-sequence activities of the parent struct
func (b *Build) Cleanup() (err error) {
	// No action occurs
	return nil
}
