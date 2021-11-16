package runtime

import (
	"io/ioutil"
	"strings"
	"udap/config"
	"udap/types"
)

// discoverModules locates and initializes modules in the plugin directory
func (r *Runtime) discoverModules() {
	dir, err := ioutil.ReadDir("./plugins")
	if err != nil {
		return
	}
	for _, info := range dir {
		if strings.Contains(info.Name(), ".so") {
			mod := types.Module{}
			mod.Path = info.Name()
			err = mod.Emplace()
			if err != nil {
				config.Error(err.Error())
			}
			r.modules[mod.Id.String()] = &mod
		}
	}
}
