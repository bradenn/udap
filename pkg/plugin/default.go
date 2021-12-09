// Copyright (c) 2021 Braden Nicholson

package plugin

type Config struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // Module, Daemon, etc.
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
}

type Module struct {
	Config
	eventHandler   *chan Event
	requestHandler chan Request
	loaded         bool
}

// Connect is called once at the launch of the module
func (m *Module) Connect(e *chan Event) (chan Request, error) {
	m.eventHandler = e
	return m.requestHandler, nil
}

func NewModule(target *Module, config Config) {
	target.Config = config
	target.loaded = true
}

func (m *Module) RegisterEntity(entity interface{}) {
	*m.eventHandler <- Event{
		Type:      "entity",
		Operation: "register",
		Body:      entity,
	}
	// cache.WatchFn(path, handleEntity)
}
