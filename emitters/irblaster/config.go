package irblaster

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/config/types"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

type Config struct {
	types.Emitter `yaml:",inline"`
	Topic         string
}

func (cfg *Config) New(c mqtt.Client) (emitters.Emitter, error) {
	return NewMQTTIRBlaster(c, cfg.Topic), nil
}

func (cfg *Config) Id() string {
	return cfg.Emitter.Id
}
