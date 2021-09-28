package main

import (
	"fmt"
	"gorm.io/gorm"
	"plugin"
	"udap/logger"
	"udap/template"
)

type Module struct {
	Persistent
	// Metadata holds generic information about the module.
	template.Metadata
	// Environment holds persistent JSON environment information, api keys, ports, etc.
	Environment string `json:"environment"`
	// Path refers to the literal name of the module
	Path string `json:"path"`
}

func (m *Module) CreateInstance(database *gorm.DB, name string) (Instance, error) {
	instance := Instance{
		Name:     name,
		Module:   *m,
		ModuleId: m.Id,
	}
	err := database.Model(&Instance{}).Create(&instance).Error
	if err != nil {
		return instance, err
	}
	return instance, nil
}

// Module interface operations

func (m *Module) Load(database *gorm.DB) error {
	module, err := m.load()
	if err != nil {
		return err
	}
	// Populate the plugin's metadata field
	m.Metadata = module.Metadata()
	// Find the module if it exists in the database
	err = database.Model(&Module{}).Where("path = ?", m.Path).First(&m).Error
	if err != nil {
		// If the module does not exist, create it
		if err == gorm.ErrRecordNotFound {
			err = m.create(database)
			// Report errors regarding the creation
			if err != nil {
				return err
			}
		}
		// Return any unhandled errors to be reported
		return err
	}
	// If the module has been loaded, or created announce it
	logger.Info("Loaded %s %s", m.Name, m.Version)
	return nil
}

// create inserts the current module into the database
func (m *Module) create(database *gorm.DB) error {
	err := database.Create(m).Error
	// Report internal errors for later diagnostic
	if err != nil {
		return fmt.Errorf("failed to create module")
	}
	return nil
}

// Initialize loads the module under an instance
func (m *Module) Initialize() (template.Mod, error) {
	module, err := m.load()
	if err != nil {
		return nil, err
	}
	module.Initialize(m.Environment)
	return module, nil
}

func (m *Module) load() (template.Mod, error) {
	filename := fmt.Sprintf(pluginFmt, m.Path, m.Path)
	p, err := plugin.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("invalid or out-dated plugin '%s'", m.Path)
	}

	lookup, err := p.Lookup("Export")
	if err != nil {
		return nil, fmt.Errorf("plugin '%s' provides not exported accesspoint", m.Path)
	}
	return lookup.(template.Mod), nil
}

func (m *Module) valid() bool {
	_, err := m.load()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// Hooks

func (m *Module) BeforeCreate(_ *gorm.DB) error {
	if !m.valid() {
		return fmt.Errorf("module does not exist")
	}
	module, err := m.load()
	if err != nil {
		return err
	}
	m.Metadata = module.Metadata()
	return nil
}

func (m *Module) AfterFind(_ *gorm.DB) error {

	return nil
}

// API

// func RouteModules(router chi.Router) {
// 	router.Post("/", createModule)
// 	router.Get("/{id}", findModule)
// 	router.Get("/", findModules)
// }
//
// func findModule(w http.ResponseWriter, r *http.Request) {
// 	req, db := server.NewRequest(w, r)
//
// 	id := req.Param("id")
//
// 	var model Module
//
// 	db.Where("id = ?", id).First(&model)
//
// 	if err := db.Error; err != nil {
// 		req.Reject(err.Error(), http.StatusNotFound)
// 		return
// 	}
//
// 	req.Resolve(model, http.StatusOK)
// }
//
// func createModule(w http.ResponseWriter, r *http.Request) {
//
// 	req, db := server.NewRequest(w, r)
//
// 	var model Module
//
// 	req.DecodeModel(&model)
// 	db.Create(&model)
//
// 	if err := db.Error; err != nil {
// 		req.Reject(err.Error(), http.StatusNotFound)
// 		return
// 	}
//
// 	req.Resolve(model, http.StatusOK)
// }
//
// func findModules(w http.ResponseWriter, r *http.Request) {
// 	req, db := server.NewRequest(w, r)
// 	var models []Module
//
// 	db.Find(&models)
//
// 	if err := db.Error; err != nil {
// 		req.Reject(err.Error(), http.StatusNotFound)
// 		return
// 	}
//
// 	req.Resolve(models, http.StatusOK)
// }
