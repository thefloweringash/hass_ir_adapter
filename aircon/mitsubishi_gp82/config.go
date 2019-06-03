package mitsubishi_gp82

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/aircon"
	"github.com/thefloweringash/hass_ir_adapter/config/types"
	"github.com/thefloweringash/hass_ir_adapter/device"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

type Config struct {
	types.Aircon `yaml:",inline"`
}

func (cfg *Config) New(c mqtt.Client, logger *log.Logger, stateDir string, emitter emitters.Emitter) (device.Device, error) {
	impl := &Device{}
	return aircon.New(&cfg.Aircon, c, emitter, impl, stateDir)
}

func (cfg *Config) Id() string {
	return cfg.Aircon.Id
}

func (cfg *Config) EmitterId() string {
	return cfg.Aircon.Emitter
}
