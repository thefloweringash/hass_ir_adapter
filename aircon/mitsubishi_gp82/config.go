package mitsubishi_gp82

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/aircon"
	"github.com/thefloweringash/hass_ir_adapter/config/types"
	"github.com/thefloweringash/hass_ir_adapter/device"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

type Config struct {
	types.Aircon `yaml:",inline"`
}

func (cfg *Config) New(
	c mqtt.Client,
	emitter emitters.Emitter,
	stateDir string,
) (device.Device, error) {
	impl := &Device{emitter}
	return aircon.New(&cfg.Aircon, c, impl, stateDir)
}

func (cfg *Config) Id() string {
	return cfg.Aircon.Id
}

func (cfg *Config) EmitterId() string {
	return cfg.Aircon.Emitter
}
