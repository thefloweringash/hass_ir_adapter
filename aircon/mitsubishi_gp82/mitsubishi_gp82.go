package mitsubishi_gp82

import (
	"fmt"

	"github.com/thefloweringash/hass_ir_adapter/aircon"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

type MitsubishiAircon struct {
	emitters.Emitter
}

type Mode byte

const (
	ModeHeating Mode = 1
	ModeDry     Mode = 2
	ModeCooling Mode = 3
)

type DryIntensity byte

const (
	DryStandard DryIntensity = 0
	DryWeak     DryIntensity = 1
	DryStrong   DryIntensity = 3
)

type WindSpeed byte

const (
	WindAuto   WindSpeed = 0
	WindQuiet  WindSpeed = 2
	WindWeak   WindSpeed = 3
	WindStrong WindSpeed = 5
)

func checksum(d []byte) byte {
	var checksum byte
	for _, b := range d {
		checksum += b
	}
	return checksum
}

func encode(state aircon.State) ([]byte, error) {
	// There's an intermediate state here that the aircon.State expands into.
	// It's not realised as a type, just as local variables.

	var timerMode byte = 0
	var timerValue byte = 0
	var isTimerCommand byte = 0
	var on byte
	var dryIntensity DryIntensity = DryStandard
	var mode Mode
	var windDirection byte = 0 // auto
	var windSpeed WindSpeed = WindAuto
	var coolFeeling byte = 0

	if state.Mode == "off" {
		on = 0
	} else {
		on = 1

	}

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

	byte1 := 1<<5 | timerMode<<3 | on<<2
	byte2 := byte(dryIntensity)<<2 | byte(mode)
	byte3 := 31 - byte(state.Temperature)
	byte4 := isTimerCommand<<6 | windDirection<<3 | byte(windSpeed)
	byte5 := coolFeeling << 5

	encoded := [14]byte{
		0x23, 0xcb, 0x26, 0x01, 0x00,
		byte1, byte2, byte3, byte4, timerValue, 0x00, byte5, 0x00,
	}
	encoded[13] = checksum(encoded[0:13])

	return encoded[:], nil
}

func (ac *MitsubishiAircon) PushState(state aircon.State) error {
	payload, err := encode(state)
	if err != nil {
		return err
	}

	cmd := emitters.Command{
		Encoding: emitters.Panasonic,
		Payload:  payload,
	}
	return ac.Emit(cmd)
}

func (ac *MitsubishiAircon) ValidState(state aircon.State) bool {
	_, err := encode(state)
	return err != nil
}

func (ac *MitsubishiAircon) DefaultState() aircon.State {
	return aircon.State{
		"off", "auto", 18,
	}
}

func (ac *MitsubishiAircon) ExportConfig(config *aircon.AirconConfig) {
	config.FanModes = []string{"auto", "quiet", "weak", "strong"}
	config.Modes = []string{"off", "cool", "heat", "dry"}
	config.MinTemp = 18
	config.MaxTemp = 31
}

func NewAircon(emitter emitters.Emitter) aircon.AirconController {
	return &MitsubishiAircon{emitter}
}
