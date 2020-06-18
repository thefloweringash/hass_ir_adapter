package tasmota_hvac

import (
	"errors"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/aircon"
	"github.com/thefloweringash/hass_ir_adapter/config/types"
	"github.com/thefloweringash/hass_ir_adapter/device"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

type Config struct {
	types.Aircon `yaml:",inline"`
	Vendor       string
	Modes        []string
	FanModes     []string
	MinTemp      float32
	MaxTemp      float32
}

func (cfg *Config) New(c mqtt.Client, logger *log.Logger, stateDir string, emitter emitters.Emitter) (device.Device, error) {
	if cfg.Vendor == "" {
		return nil, errors.New("tasmota_hvac missing vendor")
	}

	modes := cfg.Modes
	if len(modes) == 0 {
		modes = []string{"off", "auto", "cool", "heat", "dry", "fan_only"}
	}

	fanModes := cfg.FanModes
	if len(fanModes) == 0 {
		fanModes = []string{"auto", "low", "med", "high"}
	}

	maxTemp := cfg.MaxTemp
	if maxTemp == 0 {
		maxTemp = 32
	}
	minTemp := cfg.MinTemp
	if minTemp == 0 {
		minTemp = 16
	}

	impl := &Device{
		Vendor:   cfg.Vendor,
		Modes:    modes,
		FanModes: fanModes,
		MinTemp:  minTemp,
		MaxTemp:  maxTemp,
	}
	return aircon.New(&cfg.Aircon, c, emitter, impl, stateDir)
}

func (cfg *Config) Id() string {
	return cfg.Aircon.Id
}

func (cfg *Config) EmitterId() string {
	return cfg.Aircon.Emitter
}
