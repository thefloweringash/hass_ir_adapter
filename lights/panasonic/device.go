package panasonic

import (
	"github.com/thefloweringash/hass_ir_adapter/device"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/emitters/encodings"
	"github.com/thefloweringash/hass_ir_adapter/lights"
)

type Device struct {
	channel Channel
}

type State struct {
	lights.State
	Brightness uint8 `json:"brightness" hass:"brightness"`
}

func (state State) Bindings() []device.Binding {
	bindings := state.State.Bindings()
	options := device.AutomaticBindingOptions{TemplateSuffix: "value"}
	bindings = append(bindings, device.AutomaticBindings(state, options)...)
	return bindings
}

func (device *Device) Config() map[string]interface{} {
	return map[string]interface{}{}
}

func (device *Device) DefaultState() device.State {
	return State{Brightness: 255}
}

func (device *Device) PushState(emitter emitters.Emitter, rawState device.State) error {
	state := rawState.(State)

	command := Command{Channel: device.channel}
	switch {
	case !state.On:
		command.Action = ActionOff
	case state.Brightness < 26:
		command.Action = ActionNight
	case state.Brightness >= 230:
		command.Action = ActionFull
	default:
		command.Action = ActionOn
	}

	encoded, err := command.Encode()
	if err != nil {
		return err
	}

	return emitter.Emit(encodings.Panasonic{
		Payload: encoded,
	})
}
