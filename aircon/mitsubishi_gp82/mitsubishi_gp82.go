package mitsubishi_gp82

import (
	"encoding/binary"
	"fmt"

	"github.com/thefloweringash/hass_ir_adapter/aircon"
	"github.com/thefloweringash/hass_ir_adapter/emitters"

	"gopkg.in/restruct.v1"
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

type FullState struct {
	On           bool
	Mode         Mode
	WindSpeed    WindSpeed
	DryIntensity DryIntensity
	Temperature  byte
}

type Packet struct {
	Header []byte `struct:"[5]byte"`

	Padding1  byte `struct:"uint8:3"`
	TimerMode byte `struct:"uint8:2"`
	OnOff     byte `struct:"uint8:1"`
	Padding2  byte `struct:"uint8:2"`

	Padding3     byte `struct:"uint8:4"`
	DryIntensity byte `struct:"uint8:2"`
	Mode         byte `struct:"uint8:2"`

	Padding4    byte `struct:"uint8:4"`
	Temperature byte `struct:"uint8:4"`

	IsTimerCommand byte `struct:"uint8:2"`
	WindDirection  byte `struct:"uint8:3"`
	WindSpeed      byte `struct:"uint8:3"`

	TimerValue byte `struct:"uint8"`

	Padding5 byte `struct:"uint8"`

	Padding6    byte `struct:"uint8:2"`
	CoolFeeling byte `struct:"uint8:1"`
	Padding7    byte `struct:"uint8:5"`

	Padding8 byte `struct:"uint8"`
}

func (config FullState) Encode() ([]byte, error) {
	var onOff byte
	if config.On {
		onOff = 1
	}

	packet := Packet{
		Header:   []byte{0x23, 0xcb, 0x26, 0x01, 0x00},
		Padding1: 1,

		OnOff:        onOff,
		Temperature:  31 - config.Temperature,
		Mode:         byte(config.Mode),
		WindSpeed:    byte(config.WindSpeed),
		DryIntensity: byte(config.DryIntensity),
	}

	packed, err := restruct.Pack(binary.LittleEndian, &packet)
	if err != nil {
		return packed, err
	}

	return append(packed, checksum(packed)), nil
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

func (ac *MitsubishiAircon) PushState(state aircon.State) error {
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

func (ac *MitsubishiAircon) ValidState(state aircon.State) bool {
	_, err := Encode(state)
	return err != nil
}

func (ac *MitsubishiAircon) DefaultState() aircon.State {
	return aircon.State{
		Mode: "off", FanMode: "auto", Temperature: 18,
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
