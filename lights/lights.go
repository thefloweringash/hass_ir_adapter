package lights

import (
	"reflect"

	"github.com/thefloweringash/hass_ir_adapter/device"
)

const (
	ColorTempBluest  = 153
	ColorTempReddest = 500
	ColorTempWhitest = (ColorTempReddest + ColorTempBluest) / 2
	BrightnessMin    = 0
	BrightnessMax    = 255
)

type State struct {
	On bool `json:"on,string"`
}

func SwitchBinding() device.Binding {
	return &device.CallbackBinding{
		Topic: "switch",
		Conf: map[string]string{
			"state_value_template": "{{ value_json.on }}",
			"payload_on":           "true",
			"payload_off":          "false",
			"command_topic":        "~/switch",
		},
		Callback: func(state device.State, value string) (device.State, error) {
			stateReflection := reflect.TypeOf(state)
			newStatePtr := reflect.New(stateReflection)
			newStatePtr.Elem().Set(reflect.ValueOf(state))

			target := newStatePtr.Elem().FieldByName("On")
			target.SetBool(value == "true")

			return newStatePtr.Elem().Interface().(device.State), nil
		},
	}
}

func (state State) Bindings() []device.Binding {
	return []device.Binding{SwitchBinding()}
}

func Lerp(a, b float64, proportion float64) float64 {
	return b*proportion + a*(1-proportion)
}

func Proportion(a, b float64, val float64) float64 {
	return (val - a) / (b - a)
}

func ColorTempProportion(val uint16) float64 {
	return Proportion(ColorTempBluest, ColorTempReddest, float64(val))
}

func BrightnessProportion(val uint8) float64 {
	return Proportion(BrightnessMin, BrightnessMax, float64(val))
}
