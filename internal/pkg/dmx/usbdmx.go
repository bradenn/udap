package dmx

import (
	"flag"
	"fmt"
	"github.com/google/gousb"
)

// Controller Generic interface for all USB DMX controllers
type Controller interface {
	Connect() (err error)
	Close() error
	GetSerial() (info string, err error)
	GetProduct() (info string, err error)
	SetChannel(channel int16, value byte) error
	GetChannel(channel int16) (byte, error)
	Render() error
}

// ControllerConfig configuration for controlling device
type ControllerConfig struct {
	OutputId   int `toml:"outputId"`
	Context    *gousb.Context
	DebugLevel int
}

// NewConfig helper function for creating a new ControllerConfig
func NewConfig(outputId int) ControllerConfig {
	outputInterfaceID := flag.Int("output-id", outputId, "Output interface ID for device")
	debugLevel := flag.Int("debug", 0, "Debug level for USB context")
	flag.Parse()
	return ControllerConfig{
		OutputId:   *outputInterfaceID,
		DebugLevel: *debugLevel,
	}
}

// ValidateDMXChannel helper function for ensuring channel is within range
func ValidateDMXChannel(channel int) (err error) {
	if channel < 1 || channel > 512 {
		return fmt.Errorf("Channel %d out of range, must be between 1 and 512", channel)
	}

	return nil
}

// GetUSBContext gets a gousb/context for a given configuration
func (c *ControllerConfig) GetUSBContext() {
	c.Context = gousb.NewContext()
	c.Context.Debug(c.DebugLevel)
}
