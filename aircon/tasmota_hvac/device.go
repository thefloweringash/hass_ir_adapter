package tasmota_hvac

import (
	"errors"
	"fmt"

	"github.com/thefloweringash/hass_ir_adapter/aircon"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/emitters/tasmota"
)

type Device struct {
	Vendor   string
	Modes    []string
	FanModes []string
	MinTemp  float32
	MaxTemp  float32
}

func (ac *Device) PushState(emitter emitters.Emitter, state aircon.State) error {
	tasmotaEmitter, ok := emitter.(*tasmota.TasmotaEmitter)

	if !ok {
		return errors.New("tasmota may only be used with a tasmota emitter")
	}

	command := tasmota.TasmotaHvac{
		Vendor:  ac.Vendor,
		Power:   tasmota.StringBool(state.Mode != "off"),
		Celsius: true,
		Temp:    state.Temperature,
	}

	switch state.Mode {
	case "off":
		command.Mode = tasmota.ModeOff
	case "auto":
		command.Mode = tasmota.ModeAuto
	case "cool":
		command.Mode = tasmota.ModeCool
	case "heat":
		command.Mode = tasmota.ModeHeat
	case "dry":
		command.Mode = tasmota.ModeDry
	case "fan_only":
		command.Mode = tasmota.ModeFan
	default:
		return fmt.Errorf("unknown mode: '%s'", state.Mode)
	}

	switch state.FanMode {
	case "auto":
		command.FanSpeed = tasmota.FanSpeedAuto
	case "low":
		command.FanSpeed = tasmota.FanSpeedLow
	case "med":
		command.FanSpeed = tasmota.FanSpeedMid
	case "high":
		command.FanSpeed = tasmota.FanSpeedHighest
	}

	return tasmotaEmitter.EmitHVAC(command)
}

func (ac *Device) ValidState(state aircon.State) bool {
	return true
}

func (ac *Device) DefaultState() aircon.State {
	return aircon.State{
		Mode: "off", FanMode: "auto", Temperature: 18,
	}
}

func (ac *Device) Config() map[string]interface{} {
	return map[string]interface{}{
		aircon.KeyModes:    ac.Modes,
		aircon.KeyFanModes: ac.FanModes,
		aircon.KeyMinTemp:  ac.MinTemp,
		aircon.KeyMaxTemp:  ac.MaxTemp,
	}
}
