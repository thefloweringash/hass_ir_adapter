package aircon

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/config/types"
	"github.com/thefloweringash/hass_ir_adapter/device"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

const (
	KeyFanModes = "fan_modes"
	KeyModes    = "modes"
	KeyMaxTemp  = "max_temp"
	KeyMinTemp  = "min_temp"
)

type State struct {
	Mode        string  `json:"mode" hass:"mode"`
	FanMode     string  `json:"fan_mode" hass:"fan_mode"`
	Temperature float32 `json:"temperature" hass:"temperature"`
}

func (state State) Bindings() []device.Binding {
	var options device.AutomaticBindingOptions
	return device.AutomaticBindings(state, options)
}

type AirconController interface {
	PushState(emitter emitters.Emitter, state State) error
	ValidState(state State) bool
	DefaultState() State
	Config() map[string]interface{}
}

type Aircon struct {
	controller       AirconController
	temperatureTopic string
}

func (aircon *Aircon) Config() map[string]interface{} {
	config := map[string]interface{}{}
	if aircon.temperatureTopic != "" {
		config["current_temperature_topic"] = aircon.temperatureTopic
	}
	for k, v := range aircon.controller.Config() {
		config[k] = v
	}
	return config
}

func (aircon *Aircon) DefaultState() device.State {
	return aircon.controller.DefaultState()
}

func (aircon *Aircon) PushState(emitter emitters.Emitter, state device.State) error {
	return aircon.controller.PushState(emitter, state.(State))
}

func New(
	cfg *types.Aircon,
	c mqtt.Client,
	emitter emitters.Emitter,
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
		emitter,
		aircon,
		stateDir,
	)
}
