package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MQTT     MQTT `yaml:",omitempty"`
	Emitters []Emitter
	Aircons  []Aircon
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

type Emitter struct {
	Id    string
	Type  string
	Topic string
}

type Aircon struct {
	Id               string
	Name             string
	Type             string
	Emitter          string
	TemperatureTopic string `yaml:"temperature_topic"`
}

func LoadConfig(configFile string) (*Config, error) {
	fileContents, err := ioutil.ReadFile(configFile)

	if err != nil {
		return nil, err
	}

	config := Config{}

	if err := yaml.UnmarshalStrict(fileContents, &config); err != nil {
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
