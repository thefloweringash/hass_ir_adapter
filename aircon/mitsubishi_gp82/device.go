package mitsubishi_gp82

import (
	"fmt"

	"github.com/thefloweringash/hass_ir_adapter/aircon"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

type Device struct {
	emitters.Emitter
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

func (ac *Device) PushState(state aircon.State) error {
	payload, err := Encode(state)
	if err != nil {
		return err
	}

	cmd := emitters.Command{
		Encoding: emitters.Panasonic,
		Payload:  payload,
	}
	return ac.Emit(cmd)
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

func (ac *Device) ExportConfig(config *aircon.Config) {
	config.FanModes = []string{"auto", "quiet", "weak", "strong"}
	config.Modes = []string{"off", "cool", "heat", "dry"}
	config.MinTemp = 18
	config.MaxTemp = 31
}
