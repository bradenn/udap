// Copyright (c) 2021 Braden Nicholson

package runtime

import (
	"io/ioutil"
	"strings"
	"udap/module"
	"udap/udap"
)

// discoverModules locates and initializes modules in the plugin directory
func (r *Runtime) discoverModules() {

	dir, err := ioutil.ReadDir("./plugins")
	if err != nil {
		return
	}
	for _, info := range dir {
		if strings.Contains(info.Name(), ".so") {
			mod := module.Module{}
			mod.Path = info.Name()
			err = mod.Emplace()
			if err != nil {
				udap.Error(err.Error())
			}
			r.modules[mod.Id] = &mod
		}
	}
}
