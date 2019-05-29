package aircon

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/config/types"
	"github.com/thefloweringash/hass_ir_adapter/device"
)

type Config struct {
	Modes    []string `json:"modes"`
	FanModes []string `json:"fan_modes"`
	MinTemp  float32  `json:"min_temp"`
	MaxTemp  float32  `json:"max_temp"`

	CurrentTemperatureTopic string `json:"current_temperature_topic,omitempty"`
}

type State struct {
	Mode        string  `json:"mode" hass:"mode"`
	FanMode     string  `json:"fan_mode" hass:"fan_mode"`
	Temperature float32 `json:"temperature" hass:"temperature"`
}

func (state State) Bindings() []device.Binding {
	return device.AutomaticBindings(state)
}

type AirconController interface {
	PushState(state State) error
	ValidState(state State) bool
	DefaultState() State
	ExportConfig(config *Config)
}

type Aircon struct {
	controller       AirconController
	temperatureTopic string
}

func (aircon *Aircon) Config() interface{} {
	config := Config{
		CurrentTemperatureTopic: aircon.temperatureTopic,
	}
	aircon.controller.ExportConfig(&config)
	return config
}

func (aircon *Aircon) DefaultState() device.State {
	return aircon.controller.DefaultState()
}

func (aircon *Aircon) PushState(state device.State) error {
	return aircon.controller.PushState(state.(State))
}

func New(
	cfg *types.Aircon,
	c mqtt.Client,
	controller AirconController,
	stateDir string,
) (device.Device, error) {
	aircon := &Aircon{
		controller:       controller,
		temperatureTopic: cfg.TemperatureTopic,
	}

	return device.New(
		c,
		cfg.Id,
		cfg.Name,
		"climate",
		aircon,
		stateDir,
	)
}
