## Module Coding Conventions

#### Handle deferred function errors

Use the built-in module error function `module.Err(err error)` for inline error handling

```go
defer m.Err(file.Close())
```

#### Example Module

Below is a fictitious module called **WifiLights**, which provides an interface for devices controlled over the network.

**Declare Entrypoint**

Udap will attempt to access a variable named `Module` which has the super-class `plugin.Module`

```go
var Module WifiLights

type WifiLights struct {
    plugin.Module
}
```

**Initialize Entrypoint**

```go
func init() {
	// Define configuration variables
    configVariables := []plugin.Variable{
        {
        Name:        "apiKey",
        Default:     "<api_key>",
        Description: "The WifiLight api access token provided from the WifiLight website.",
        },
    }
	
    // Define module configuration
    Module.Config = plugin.Config{
        Name:        "wifi-lights",
        Type:        "module",
        Description: "Control light over wifi",
        Version:     "0.0.1",
        Author:      "<author>",
        Variables:   configVariables,
    }
}
```

**Define Module Overrides**

```go

func (w *WifiLights) Setup() error {
	// Do one-time setup tasks

	// Set the update interval to update every 2 seconds
	err := v.UpdateInterval(2000)
	if err != nil {
        return err
    }
	return nil
}

func (w *WifiLights) Update() error {
    if v.Ready() {
        // Push changes to attributes
    }
    return nil
}

func (w *WifiLights) Run() error {
    // Run start up sequences, initialize entities, attributes, etc
    return nil
}

func (w *WifiLights) Dispose() error {
    // Shutdown any components and close files if needed
    return nil
}
```