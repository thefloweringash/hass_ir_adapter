package tasmota

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/config/types"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

type Config struct {
	types.Emitter `yaml:",inline"`
	Topic         string
}

func (cfg *Config) New(c mqtt.Client, logger *log.Logger) (emitters.Emitter, error) {
	return NewTasmotaEmitter(c, logger, cfg.Topic), nil
}

func (cfg *Config) Id() string {
	return cfg.Emitter.Id
}
