package daikin

import (
	"fmt"

	"github.com/thefloweringash/hass_ir_adapter/aircon"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/emitters/encodings"
)

type Device struct {
	emitters.Emitter
}

func encode(state aircon.State) ([]byte, []byte, error) {
	var mode Mode

	switch state.Mode {
	case "off", "auto":
		// always set a mode to be a valid packet
		mode = ModeAuto
	case "dry":
		mode = ModeDry
	case "cool":
		mode = ModeCooling
	case "heat":
		mode = ModeHeating
	case "sendoff_wind":
		mode = ModeSendoffWind
	default:
		return nil, nil, fmt.Errorf("unknown mode: '%s'", state.Mode)
	}

	return State{
		On:          state.Mode != "off",
		Mode:        mode,
		Temperature: state.Temperature,
	}.Encode()
}

func (ac *Device) Config() map[string]interface{} {
	return map[string]interface{}{
		aircon.KeyMinTemp: 18,
		aircon.KeyMaxTemp: 31,
		aircon.KeyModes:   []string{"off", "dry", "cool", "heat"},
	}
}

func (ac *Device) DefaultState() aircon.State {
	return aircon.State{
		Mode:        "off",
		Temperature: 18,
	}
}

func (ac *Device) ValidState(state aircon.State) bool {
	// TODO
	return true
}

func (ac *Device) PushState(emitter emitters.Emitter, state aircon.State) error {
	p1, p2, err := encode(state)
	if err != nil {
		return err
	}

	command := encodings.Daikin{P1: p1, P2: p2}
	return emitter.Emit(command)
}
