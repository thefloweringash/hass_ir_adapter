package mitsubishi_gp82

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
	var windSpeed WindSpeed

	switch state.Mode {
	case "off", "heat":
		// always set a mode to be a valid packet
		mode = ModeHeating
	case "dry":
		mode = ModeDry
	case "cool":
		mode = ModeCooling
	default:
		return nil, fmt.Errorf("Unknown mode: '%s'", state.Mode)
	}

	switch state.FanMode {
	case "auto":
		windSpeed = WindAuto
	case "quiet":
		windSpeed = WindQuiet
	case "weak":
		windSpeed = WindWeak
	case "strong":
		windSpeed = WindStrong
	default:
		return nil, fmt.Errorf("Unknown fan mode: '%s'", state.FanMode)
	}

	return FullState{
		On:          state.Mode != "off",
		Mode:        mode,
		Temperature: byte(state.Temperature),
		WindSpeed:   windSpeed,
	}.Encode()
}

func (ac *Device) PushState(emitter emitters.Emitter, state aircon.State) error {
	payload, err := Encode(state)
	if err != nil {
		return err
	}

	return emitter.Emit(encodings.Panasonic{
		Payload: payload,
	})
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
		aircon.KeyFanModes: []string{"auto", "quiet", "weak", "strong"},
		aircon.KeyModes:    []string{"off", "cool", "heat", "dry"},
		aircon.KeyMinTemp:  18,
		aircon.KeyMaxTemp:  31,
	}
}
