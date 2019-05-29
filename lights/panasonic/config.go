package panasonic

import (
	"fmt"

	"github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v3"

	"github.com/thefloweringash/hass_ir_adapter/config/types"
	"github.com/thefloweringash/hass_ir_adapter/device"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/lights"
)

type Config struct {
	types.Light `yaml:",inline"`
	Channel     Channel
}

type rawConfig struct {
	types.Light `yaml:",inline"`
	Channel     int
}

func (config *Config) UnmarshalYAML(value *yaml.Node) error {
	raw := rawConfig{}
	if err := value.Decode(&raw); err != nil {
		return err
	}
	config.Light = raw.Light
	switch raw.Channel {
	case 1:
		config.Channel = Channel1
	case 2:
		config.Channel = Channel2
	case 3:
		config.Channel = Channel3
	default:
		return fmt.Errorf("invalid channel value: %d", raw.Channel)
	}
	return nil
}

func (config *Config) New(c mqtt.Client, emitter emitters.Emitter, stateDir string) (device.Device, error) {
	impl := &Device{emitter, config.Channel}
	return lights.New(&config.Light, c, impl, stateDir)
}

func (config *Config) Id() string {
	return config.Light.Id
}

func (config *Config) EmitterId() string {
	return config.Light.Emitter
}
