// Copyright (c) 2021 Braden Nicholson

package plugin

import (
	"fmt"
	"io/ioutil"
	"plugin"
	"strings"
)

// Find attempts to discover plugins in
func Find(path string) (plugins map[string]*Plugin, err error) {
	plugins = map[string]*Plugin{}
	// Attempt to read the path directory
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	// Look for a file ending with .so
	for _, info := range dir {
		// If the info of this file has a suffix of .so, we load the plugin
		if info.IsDir() {
			nDir, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", path, info.Name()))
			if err != nil {
				return plugins, err
			}
			for _, fileInfo := range nDir {
				if strings.Contains(fileInfo.Name(), ".so") {
					// load will attempt to load the plugin from the path
					name := strings.Replace(fileInfo.Name(), ".so", "", 1)
					p, err := load(fileInfo.Name(), fmt.Sprintf("%s/%s", path, info.Name()))
					if err != nil {
						return nil, err
					}
					plugins[fmt.Sprintf("%s/%s/%s", path, name, fileInfo.Name())] = &p
				}
			}
		}
	}
	// Return plugins no errors
	return plugins, nil
}

// getPlugin attempts to load the plugin from a given path
func load(name string, path string) (pl Plugin, err error) {
	// Attempt to open that plugin
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open plugin at path '%s': %s", path, err.Error())
	}
	// Attempt to access the Plugin variable to interface with the code
	lookup, err := p.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("plugin '%s' does not define a Plugin interface", path)
	}
	pl = lookup.(Plugin)
	// Return no errors
	return lookup.(Plugin), nil
}
