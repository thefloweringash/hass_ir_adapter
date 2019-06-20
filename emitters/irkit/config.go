package irkit

import (
	"log"

	"github.com/thefloweringash/hass_ir_adapter/config/types"
	"github.com/thefloweringash/hass_ir_adapter/emitters"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
	types.Emitter `yaml:",inline"`
	Endpoint      string
}

func (cfg *Config) New(c mqtt.Client, logger *log.Logger) (emitters.Emitter, error) {
	return New(cfg.Endpoint, logger), nil
}

func (cfg *Config) Id() string {
	return cfg.Emitter.Id
}
