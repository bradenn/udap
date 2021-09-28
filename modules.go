package main

import (
	"gorm.io/gorm"
	"io/ioutil"
	"udap/logger"
)

const (
	pathFmt   = "./plugins"
	pluginFmt = pathFmt + "/%s/%s.so"
)

func DiscoverModules(database *gorm.DB) {
	dir, err := ioutil.ReadDir(pathFmt)
	if err != nil {
		return
	}

	for _, info := range dir {
		if info.IsDir() {
			mod := Module{}
			mod.Path = info.Name()
			err = mod.Load(database)
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}
}
