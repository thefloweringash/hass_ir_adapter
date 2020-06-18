package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/thefloweringash/hass_ir_adapter/aircon/daikin"
	"github.com/thefloweringash/hass_ir_adapter/aircon/mitsubishi_gp82"
	"github.com/thefloweringash/hass_ir_adapter/aircon/mitsubishi_rh101"
	"github.com/thefloweringash/hass_ir_adapter/aircon/tasmota_hvac"
	"github.com/thefloweringash/hass_ir_adapter/device"
	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/emitters/irblaster"
	"github.com/thefloweringash/hass_ir_adapter/emitters/irkit"
	"github.com/thefloweringash/hass_ir_adapter/emitters/tasmota"
	"github.com/thefloweringash/hass_ir_adapter/lights/daiko"
	"github.com/thefloweringash/hass_ir_adapter/lights/panasonic"
)

type Config struct {
	MQTT     MQTT `yaml:",omitempty"`
	Emitters []emitters.Factory
	Devices  []device.Factory
	Lights   []device.Factory
}

type GenericConfig struct {
	MQTT     MQTT
	Emitters []yaml.Node
	Aircons  []yaml.Node
	Lights   []yaml.Node
}

func getTypeKey(value yaml.Node) (string, error) {
	asMap := map[string]yaml.Node{}
	if err := value.Decode(&asMap); err != nil {
		return "", err
	}
	typeNode, ok := asMap["type"]
	if !ok {
		return "", fmt.Errorf("missing type entry")
	}
	typeString := ""
	if err := typeNode.Decode(&typeString); err != nil {
		return "", err
	}
	return typeString, nil
}

func decodeVirtual(value yaml.Node, newFactory func(t string) interface{}) (interface{}, error) {
	var factory interface{}

	typeString, err := getTypeKey(value)
	if err != nil {
		return factory, err
	}

	factory = newFactory(typeString)
	if factory == nil {
		return factory, fmt.Errorf("unknown type: %s", typeString)
	}

	err = value.Decode(factory)
	return factory, err
}

func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	generic := GenericConfig{}
	err := value.Decode(&generic)
	if err != nil {
		return err
	}

	c.MQTT = generic.MQTT

	for _, emitterNode := range generic.Emitters {
		factory, err := decodeVirtual(emitterNode, func(typeStr string) (factory interface{}) {
			switch typeStr {
			case "irblaster":
				factory = &irblaster.Config{}
			case "irkit":
				factory = &irkit.Config{}
			case "tasmota":
				factory = &tasmota.Config{}
			}
			return
		})

		if err != nil {
			return err
		}

		c.Emitters = append(c.Emitters, factory.(emitters.Factory))
	}

	for _, airconNode := range generic.Aircons {
		factory, err := decodeVirtual(airconNode, func(typeStr string) (factory interface{}) {
			switch typeStr {
			case "mitsubishi_rh101":
				factory = &mitsubishi_rh101.Config{}
			case "mitsubishi_gp82":
				factory = &mitsubishi_gp82.Config{}
			case "daikin":
				factory = &daikin.Config{}
			case "tasmota_hvac":
				factory = &tasmota_hvac.Config{}
			}
			return
		})

		if err != nil {
			return err
		}

		c.Devices = append(c.Devices, factory.(device.Factory))
	}

	for _, lightNode := range generic.Lights {
		factory, err := decodeVirtual(lightNode, func(typeStr string) (factory interface{}) {
			switch typeStr {
			case "panasonic":
				factory = &panasonic.Config{}
			case "daiko":
				factory = &daiko.Config{}
			}
			return
		})

		if err != nil {
			return err
		}

		c.Devices = append(c.Devices, factory.(device.Factory))
	}

	return nil
}

type MQTT struct {
	Broker   string
	Username string
	Password string
}

func (mqtt *MQTT) isZero() bool {
	return mqtt.Broker == "" &&
		mqtt.Username == "" &&
		mqtt.Password == ""

}

func LoadConfig(configFile string) (*Config, error) {
	fileContents, err := ioutil.ReadFile(configFile)

	if err != nil {
		return nil, err
	}

	config := Config{}

	if err := yaml.Unmarshal(fileContents, &config); err != nil {
		return nil, err
	}

	if config.MQTT.isZero() {
		config.MQTT = MQTT{
			Broker:   "tcp://localhost:1883",
			Username: "hass_ir_adapter",
			Password: "hass_ir_adapter",
		}
	}

	return &config, nil
}
