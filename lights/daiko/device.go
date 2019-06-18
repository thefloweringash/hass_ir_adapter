package daiko

import (
	"log"
	"math"

	"github.com/thefloweringash/hass_ir_adapter/device"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/emitters/encodings"
	"github.com/thefloweringash/hass_ir_adapter/lights"
)

type Device struct {
	logger  *log.Logger
	channel Channel
}

type State struct {
	lights.State
	Brightness uint8  `json:"brightness" hass:"brightness"`
	ColorTemp  uint16 `json:"color_temp" hass:"color_temp"`
}

func updateColorTempHook(state device.State, update func(device.State) (device.State, error)) (device.State, error) {
	updatedState, err := update(state)
	if err != nil {
		return nil, err
	}

	newState := updatedState.(State)

	brightness := BrightnessFromHass(newState.Brightness)
	warmth := WarmthFromHass(newState.ColorTemp)

	clampedBrightness := ClampBrightness(warmth, brightness)
	if brightness != clampedBrightness {
		/*
			fmt.Printf("clamped brightness: %d -> %d (%d -> %d)\n",
				warmth, clampedBrightness,
				newState.ColorTemp, BrightnessToHass(clampedBrightness),
			)
		*/

		newState.Brightness = BrightnessToHass(clampedBrightness)
	}

	return newState, nil
}

func updateBrightnessHook(state device.State, update func(device.State) (device.State, error)) (device.State, error) {
	updatedState, err := update(state)
	if err != nil {
		return nil, err
	}

	newState := updatedState.(State)

	brightness := BrightnessFromHass(newState.Brightness)
	warmth := WarmthFromHass(newState.ColorTemp)

	clampedWarmth := ClampWarmth(warmth, brightness)
	if warmth != clampedWarmth {
		/*
			fmt.Printf("clamped warmth: %d -> %d (%d -> %d)\n",
				warmth, clampedWarmth,
				newState.ColorTemp, WarmthToHass(clampedWarmth),
			)
		*/

		newState.ColorTemp = WarmthToHass(clampedWarmth)
	}

	return newState, nil
}

func (state State) Bindings() []device.Binding {
	bindings := state.State.Bindings()

	options := device.AutomaticBindingOptions{
		TemplateSuffix: "value",
		UpdateHooks: map[string]device.BindingHook{
			"Brightness": updateBrightnessHook,
			"ColorTemp":  updateColorTempHook,
		},
	}

	bindings = append(bindings, device.AutomaticBindings(state, options)...)
	return bindings
}

func (device *Device) Config() map[string]interface{} {
	return map[string]interface{}{}
}

func (device *Device) DefaultState() device.State {
	return State{
		Brightness: lights.BrightnessMax,
		ColorTemp:  lights.ColorTempWhitest,
	}
}

func WarmthFromHass(colorTemp uint16) uint8 {
	proportion := lights.ColorTempProportion(colorTemp)
	return uint8(math.Round(lights.Lerp(WarmthMin, WarmthMax, proportion)))
}

func WarmthToHass(warmth uint8) uint16 {
	proportion := lights.Proportion(WarmthMin, WarmthMax, float64(warmth))
	return uint16(math.Round(lights.Lerp(lights.ColorTempBluest, lights.ColorTempReddest, proportion)))
}

func BrightnessFromHass(brightness uint8) uint8 {
	proportion := lights.BrightnessProportion(brightness)
	return uint8(math.Round(lights.Lerp(BrightnessMin, BrightnessMax, proportion)))
}

func BrightnessToHass(brightness uint8) uint8 {
	proportion := lights.Proportion(BrightnessMin, BrightnessMax, float64(brightness))
	return uint8(math.Round(lights.Lerp(lights.BrightnessMin, lights.BrightnessMax, proportion)))
}

func (device *Device) PushState(emitter emitters.Emitter, rawState device.State) error {
	state := rawState.(State)

	var command []byte
	var err error

	if !state.On {
		command = Off(device.channel)
	} else {
		brightness := BrightnessFromHass(state.Brightness)
		warmth := WarmthFromHass(state.ColorTemp)

		device.logger.Printf("transformed brightness:%d, color_temp:%d to warmth:%d, brightness:%d",
			state.Brightness, state.ColorTemp, warmth, brightness)

		command, err = On(device.channel, warmth, brightness)
	}

	if err != nil {
		return err
	}

	return emitter.Emit(encodings.Repeat(
		encodings.NEC{Payload: command},
		35,
	))
}
