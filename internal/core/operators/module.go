// Copyright (c) 2022 Braden Nicholson

package operators

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/core/ports"
	"udap/internal/log"
	"udap/internal/plugin"
)

const PATH = "./modules"

type lifeCycle struct {
	iface plugin.ModuleInterface
	mutex sync.Mutex
}

func newLifeCycle(iface plugin.ModuleInterface) *lifeCycle {
	return &lifeCycle{
		iface: iface,
		mutex: sync.Mutex{},
	}
}

type moduleRuntime struct {
	ctrl        *controller.Controller
	runtime     map[string]*lifeCycle
	moduleMutex sync.RWMutex
}

func NewModuleOperator(ctrl *controller.Controller) ports.ModuleOperator {
	return &moduleRuntime{
		ctrl:        ctrl,
		runtime:     map[string]*lifeCycle{},
		moduleMutex: sync.RWMutex{},
	}
}

func (m *moduleRuntime) HandleEmit(mutation domain.Mutation) error {
	return nil
}

func (m *moduleRuntime) getModule(id string) (plugin.ModuleInterface, error) {

	runtime := m.runtime[id]
	if runtime == nil {
		return nil, fmt.Errorf("module lifecycle not found")
	}

	var iface plugin.ModuleInterface

	runtime.mutex.Lock()
	iface = runtime.iface
	runtime.mutex.Unlock()
	if iface == nil {

		return nil, fmt.Errorf("module not found")
	}

	return iface, nil
}

func (m *moduleRuntime) returnModule(id string, moduleInterface plugin.ModuleInterface) error {

	runtime := m.runtime[id]
	if runtime == nil {
		return fmt.Errorf("module lifecycle not found")
	}

	var iface plugin.ModuleInterface
	runtime.mutex.Lock()
	runtime.iface = moduleInterface
	runtime.mutex.Unlock()

	if iface == nil {
		return fmt.Errorf("module not found")
	}

	return nil
}

func (m *moduleRuntime) setModule(id string, moduleInterface plugin.ModuleInterface) error {
	m.runtime[id] = newLifeCycle(moduleInterface)
	return nil
}

func (m *moduleRuntime) removeModule(id string) error {
	runtime := m.runtime[id]
	if runtime == nil {
		return fmt.Errorf("module lifecycle not found")
	}
	delete(m.runtime, id)
	return nil
}

// hashFile returns a sha256 hash of a source file
func hashFile(source string) (string, error) {
	f, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}
	return string(h.Sum(nil)), nil
}

// generateSessionId provides the first bytes of the module uuid
func generateSessionId(uuid string) string {
	return strings.Split(uuid, "-")[0]
}

// generateSourcePath generates the path to a module's source code
func generateSourcePath(module string) string {
	return fmt.Sprintf("%s/%s/%s.go", PATH, module, module)
}

// generateBuildPath generates the path to a module's binary file
func generateBuildPath(module string, uuid string) string {
	return fmt.Sprintf("%s/%s/%s-%s.so", PATH, module, module, uuid)
}

// Build will compile a valid plugin file into a readable binary
func (m *moduleRuntime) Build(module string, uuid string) error {
	// Generate the source file path from the module name
	sourcePath := generateSourcePath(module)
	// Generate the output build file path
	buildPath := generateBuildPath(module, uuid)
	// Confirm that the source file exists
	if _, err := os.Stat(sourcePath); err != nil {
		return err
	}
	// Create a timeout to prevent modules from taking too long to build
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*25)
	// Cancel the timeout of it exits before the timeout is up
	defer cancelFunc()
	// Prepare the command arguments
	args := []string{"build", "-v", "-buildmode=plugin", "-o", buildPath, sourcePath}
	// Initialize the command structure
	cmd := exec.CommandContext(timeout, "go", args...)
	// Run and get the stdout and stderr from the output
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("module build failed: \n%s", out)
	}
	return nil
}

// Cleanup deletes all compiled binaries in the module folder (.so)
func (m *moduleRuntime) Cleanup(module string) error {
	// Get the path to the target directory
	directory := fmt.Sprintf("%s/%s", PATH, module)
	// Read the directory, get an array of all file entries
	dir, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	// Create an array to store the filenames
	var toDelete []string
	// Go through each entry and find out if it is a binary
	for _, entry := range dir {
		// Check if the filename ends with the extension '.so'
		if strings.HasSuffix(entry.Name(), ".so") {
			// Append to delete list with the path
			toDelete = append(toDelete, entry.Name())
		}
	}
	// Go through each of the selected files and delete them
	for _, name := range toDelete {
		// Form the path to the target file
		path := fmt.Sprintf("%s/%s/%s", PATH, module, name)
		// Remove the file
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}
	// Exit normally
	return nil

}

// Load is used to find a pre-built plugin file, and load it into the local system.
// The module reference should be saved to the repository after loading.
func (m *moduleRuntime) Load(module string, uuid string) (domain.ModuleConfig, error) {
	// Create the binary file path
	binary := generateBuildPath(module, uuid)
	// Attempt to load the plugin binary
	p, err := plugin.Load(binary)
	if err != nil {
		return domain.ModuleConfig{}, err
	}
	// Extract the plugin interface
	mod := p.(plugin.ModuleInterface)
	if mod == nil {
		return domain.ModuleConfig{}, fmt.Errorf("cannot read module")
	}
	// Connect the module to the UDAP runtime
	err = mod.Connect(m.ctrl, uuid)
	if err != nil {
		return domain.ModuleConfig{}, err
	}
	// Run the setup method
	setup, err := mod.Setup()
	if err != nil {
		return domain.ModuleConfig{}, err
	}
	// Emplace the module into the local buffer
	err = m.setModule(uuid, mod)
	if err != nil {
		return domain.ModuleConfig{}, err
	}

	marshal, err := json.Marshal(setup.Variables)
	if err != nil {
		return domain.ModuleConfig{}, err
	}

	conf := domain.ModuleConfig{
		Name:        setup.Name,
		Type:        setup.Type,
		Description: setup.Description,
		Interval:    setup.Interval,
		Version:     setup.Version,
		Author:      setup.Author,
		Variables:   string(marshal),
	}
	return conf, nil
}

// Dispose is called at the end of the lifecycle, it attempts to halt activity.
func (m *moduleRuntime) Dispose(module string, uuid string) error {
	// Get the local module
	local, err := m.getModule(uuid)
	if err != nil {
		return err
	}

	binaryPath := generateBuildPath(module, uuid)
	// Confirm that the binary file exists
	if _, err = os.Stat(binaryPath); err != nil {
		return err
	}
	// Dispose of the local module
	err = local.Dispose()
	if err != nil {
		return err
	}
	// Remove the file when the function exits
	defer func() {
		err = os.Remove(binaryPath)
		if err != nil {
			log.Err(err)
			return
		}
		err = m.Cleanup(module)
		if err != nil {
			return
		}
	}()

	err = m.removeModule(uuid)
	if err != nil {
		return err
	}

	return nil
}

// Run attempts to initialize the plugin's runtime
func (m *moduleRuntime) Run(uuid string) error {
	// Get the local module
	local, err := m.getModule(uuid)
	if err != nil {
		return err
	}
	//defer func() {
	//	err = m.returnModule(uuid, local)
	//	if err != nil {
	//		log.Err(err)
	//	}
	//}()
	// Run the local module
	err = local.Run()
	if err != nil {
		return err
	}
	return nil
}

// Update is called on every tick, it is the plugin's decision whether to use the time or defer.
func (m *moduleRuntime) Update(uuid string) error {
	// Get the local module
	local, err := m.getModule(uuid)
	if err != nil {
		return err
	}
	//defer func() {
	//	err = m.returnModule(uuid, local)
	//	if err != nil {
	//		log.Err(err)
	//	}
	//}()
	// Run the update
	err = local.Update()
	if err != nil {
		return err
	}

	return nil
}
