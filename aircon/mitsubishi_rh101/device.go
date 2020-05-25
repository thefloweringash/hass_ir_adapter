package mitsubishi_rh101

import (
	"fmt"

	"github.com/thefloweringash/hass_ir_adapter/aircon"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/emitters/encodings"
)

type Device struct {
}

func Encode(state aircon.State) ([]byte, error) {
	var mode Mode

	switch state.Mode {
	case "off", "heat":
		// always set a mode to be a valid packet
		mode = ModeHeating
	case "cool":
		mode = ModeCooling
	default:
		return nil, fmt.Errorf("Unknown mode: '%s'", state.Mode)
	}

	return FullState{
		On:          state.Mode != "off",
		Mode:        mode,
		Temperature: byte(state.Temperature),
	}.Encode()
}

func (ac *Device) PushState(emitter emitters.Emitter, state aircon.State) error {
	payload, err := Encode(state)
	if err != nil {
		return err
	}

	return emitter.Emit(encodings.Repeat(
		encodings.Panasonic{Payload: payload},
		13,
	))
}

func (ac *Device) ValidState(state aircon.State) bool {
	_, err := Encode(state)
	return err != nil
}

func (ac *Device) DefaultState() aircon.State {
	return aircon.State{
		Mode: "off", FanMode: "auto", Temperature: 18,
	}
}

func (ac *Device) Config() map[string]interface{} {
	return map[string]interface{}{
		aircon.KeyModes:   []string{"off", "cool", "heat"},
		aircon.KeyMinTemp: 18,
		aircon.KeyMaxTemp: 31,
	}
}
