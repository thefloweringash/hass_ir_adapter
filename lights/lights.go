package lights

import (
	"github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/config/types"
	"github.com/thefloweringash/hass_ir_adapter/device"
)

type Controller interface {
	PushState(state State) error
}

type Lights struct {
	controller Controller
}

type State struct {
	On         bool  `json:"on,string"`
	Brightness uint8 `json:"brightness" hass:"brightness"`
}

func (state State) Bindings() []device.Binding {
	bindings := device.AutomaticBindings(state)
	bindings = append(bindings, &device.CallbackBinding{
		Topic: "switch",
		Conf: map[string]string{
			"state_value_template": "{{ value_json.on }}",
			"payload_on":           "true",
			"payload_off":          "false",
			"command_topic":        "~/switch",
		},
		Callback: func(state device.State, value string) (device.State, error) {
			newState := state.(State)
			newState.On = value == "true"
			return newState, nil
		},
	})
	return bindings
}

type Config struct{}

func (l *Lights) Config() interface{} {
	return Config{}
}

func (l *Lights) DefaultState() device.State {
	return State{}
}

func (l *Lights) PushState(state device.State) error {
	return l.controller.PushState(state.(State))
}

func New(config *types.Light, c mqtt.Client, controller Controller, stateDir string) (device.Device, error) {
	lights := &Lights{controller: controller}
	return device.New(
		c,
		config.Id,
		config.Name,
		"light",
		lights,
		stateDir,
	)
}
