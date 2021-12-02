// Copyright (c) 2021 Braden Nicholson

package plugin

import (
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"udap/internal/log"
)

func goFile(wg *sync.WaitGroup) func(path string, fileInfo fs.DirEntry, err error) error {
	return func(path string, fileInfo fs.DirEntry, err error) error {
		if strings.HasSuffix(fileInfo.Name(), ".go") {
			// name := strings.Replace(fileInfo.Name(), ".go", "", 2)
			wg.Add(1)
			go buildFromSource(path, wg)
		}
		return nil
	}
}

func soFile(plugins *map[string]Plugin) func(path string, fileInfo fs.DirEntry, err error) error {
	return func(path string, fileInfo fs.DirEntry, err error) error {
		if strings.HasSuffix(fileInfo.Name(), ".so") {
			// load will attempt to load the plugin from the path
			name := strings.Replace(fileInfo.Name(), ".so", "", 1)
			p, err := load(name, path)
			if err != nil {
				return err
			}
			(*plugins)[path] = p
		}
		return nil
	}
}

func buildFromSource(path string, wg *sync.WaitGroup) {

	out := strings.Replace(path, ".go", ".so", 3)
	cmd := exec.Command("/Users/bradennicholson/go/go1.17.2/bin/go", "build", "-v", "-buildmode=plugin", "-o",
		out, path)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Err(err)
	}
	wg.Done()

	// log.Log("Built %s %s", name, output)
}

func BuildAll(dir string, wg *sync.WaitGroup) {
	err := filepath.WalkDir(dir, goFile(wg))
	if err != nil {
		panic(err)
	}
}

func LoadAll(dir string) (str map[string]Plugin) {
	str = map[string]Plugin{}
	err := filepath.WalkDir(dir, soFile(&str))
	if err != nil {
		panic(err)
	}
	return str
}
