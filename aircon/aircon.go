package aircon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"reflect"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

type AirconController interface {
	PushState(state State) error
	ValidState(state State) bool
	DefaultState() State
	ExportConfig(config *AirconConfig)
}

type AirconConfig struct {
	Name     string   `json:"name"`
	Modes    []string `json:"modes"`
	FanModes []string `json:"fan_modes"`
	MinTemp  int      `json:"min_temp"`
	MaxTemp  int      `json:"max_temp"`

	Prefix                   string `json:"~"`
	ModeCommandTopic         string `json:"mode_command_topic"`
	ModeStateTopic           string `json:"mode_state_topic"`
	ModeStateTemplate        string `json:"mode_state_template"`
	TemperatureCommandTopic  string `json:"temperature_command_topic"`
	TemperatureStateTopic    string `json:"temperature_state_topic"`
	TemperatureStateTemplate string `json:"temperature_state_template"`
	FanModeCommandTopic      string `json:"fan_mode_command_topic"`
	FanModeStateTopic        string `json:"fan_mode_state_topic"`
	FanModeStateTemplate     string `json:"fan_mode_state_template"`

	CurrentTemperatureTopic string `json:"current_temperature_topic"`
}

type Aircon struct {
	id               string
	name             string
	emitter          emitters.Emitter
	mqttClient       mqtt.Client
	impl             AirconController
	temperatureTopic string
	stateFile        string
}

type State struct {
	Mode        string `json:"mode"`
	FanMode     string `json:"fan_mode"`
	Temperature int    `json:"temperature"`
}

func prefix(name string) string {
	return fmt.Sprintf("homeassistant/climate/%s", name)
}

func topicName(name string, topic Topic) string {
	return fmt.Sprintf("%s/%s", prefix(name), topic)
}

func NewAircon(
	c mqtt.Client,
	emitter emitters.Emitter,
	impl AirconController,
	id string,
	name string,
	temperatureTopic string,
	stateDir string,
) (*Aircon, error) {
	return &Aircon{
		id:               id,
		name:             name,
		emitter:          emitter,
		mqttClient:       c,
		impl:             impl,
		temperatureTopic: temperatureTopic,
		stateFile:        path.Join(stateDir, id),
	}, nil
}

func (aircon *Aircon) loadState() (State, error) {
	state := aircon.impl.DefaultState()

	contents, err := ioutil.ReadFile(aircon.stateFile)

	if err != nil {
		return state, err
	}

	if err := json.Unmarshal(contents, &state); err != nil {
		return state, err
	}

	return state, nil
}

func (aircon *Aircon) writeState(state State) error {
	contents, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(aircon.stateFile, contents, 0640)
}

func (aircon *Aircon) publishState(state State) error {
	value, err := json.Marshal(state)
	if err != nil {
		return err
	}

	token := aircon.mqttClient.Publish(topicName(aircon.id, "state"), 0, true, value)
	token.Wait()
	return token.Error()
}

func stateTemplate(fieldName string) (string, error) {
	stateReflection := reflect.TypeOf(State{})
	field, ok := stateReflection.FieldByName(fieldName)
	if !ok {
		return "", fmt.Errorf("Missing state field '%v'", field)
	}
	name, ok := field.Tag.Lookup("json")
	if !ok {
		return "", fmt.Errorf("Missing json tag on state field '%v'", field)
	}
	return fmt.Sprintf("{{ value_json.%s }}", name), nil
}

type Topic string

const (
	ConfigTopic             Topic = "config"
	TemperatureCommandTopic Topic = "temperature_command"
	FanModeCommandTopic     Topic = "fan_mode_command"
	ModeCommandTopic        Topic = "mode_command"
	StateTopic              Topic = "state"
)

func (aircon *Aircon) publishConfig() error {
	modeStateTemplate, err := stateTemplate("Mode")
	if err != nil {
		return err
	}

	fanModeStateTemplate, err := stateTemplate("FanMode")
	if err != nil {
		return err
	}

	temperatureStateTemplate, err := stateTemplate("Temperature")
	if err != nil {
		return err
	}

	config := AirconConfig{
		Name:   aircon.name,
		Prefix: prefix(aircon.id),

		ModeCommandTopic:  "~/" + string(ModeCommandTopic),
		ModeStateTopic:    "~/" + string(StateTopic),
		ModeStateTemplate: modeStateTemplate,

		TemperatureCommandTopic:  "~/" + string(TemperatureCommandTopic),
		TemperatureStateTopic:    "~/" + string(StateTopic),
		TemperatureStateTemplate: temperatureStateTemplate,

		FanModeCommandTopic:  "~/" + string(FanModeCommandTopic),
		FanModeStateTopic:    "~/" + string(StateTopic),
		FanModeStateTemplate: fanModeStateTemplate,

		CurrentTemperatureTopic: aircon.temperatureTopic,
	}
	aircon.impl.ExportConfig(&config)

	value, err := json.Marshal(config)
	if err != nil {
		return err
	}

	token := aircon.mqttClient.Publish(topicName(aircon.id, ConfigTopic), 0, true, value)
	token.Wait()
	return token.Error()
}

func (aircon *Aircon) removeConfig() error {
	token := aircon.mqttClient.Publish(topicName(aircon.id, ConfigTopic), 0, true, []byte{})
	token.Wait()
	return token.Error()
}

func (aircon *Aircon) removeState() error {
	token := aircon.mqttClient.Publish(topicName(aircon.id, StateTopic), 0, true, []byte{})
	token.Wait()
	return token.Error()
}

func (aircon *Aircon) Run() (func(), error) {
	c := aircon.mqttClient
	name := aircon.id
	impl := aircon.impl

	stopChan := make(chan bool, 1)
	stopDoneChan := make(chan bool, 1)

	tempChan := make(chan int)
	setTemp := func(c mqtt.Client, m mqtt.Message) {
		temp, err := strconv.ParseFloat(string(m.Payload()), 64)
		if err != nil {
			log.Printf("Unable to understand temperature set request: %s", err)
			return
		}
		tempChan <- int(temp)
	}

	modeChan := make(chan string)
	setMode := func(c mqtt.Client, m mqtt.Message) {
		modeChan <- string(m.Payload())
	}

	fanModeChan := make(chan string)
	setFan := func(c mqtt.Client, m mqtt.Message) {
		fanModeChan <- string(m.Payload())
	}

	if token := c.Subscribe(topicName(name, TemperatureCommandTopic), 0, setTemp); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	if token := c.Subscribe(topicName(name, ModeCommandTopic), 0, setMode); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	if token := c.Subscribe(topicName(name, FanModeCommandTopic), 0, setFan); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	initialState, err := aircon.loadState()

	if err != nil {
		log.Printf("Error loading initial state: %s", err)
		log.Printf("Using default state")
	}

	if err := aircon.publishConfig(); err != nil {
		return nil, err
	}

	go func() {
		state := initialState
		running := true
		for running == true {
			if err := aircon.publishState(state); err != nil {
				log.Printf("Error publishing state: %s", err)
			}

			newState := state
			select {
			case temp := <-tempChan:
				newState.Temperature = temp
			case mode := <-modeChan:
				newState.Mode = mode
			case fanMode := <-fanModeChan:
				newState.FanMode = fanMode
			case <-stopChan:
				running = false
				continue
			}

			if err := impl.PushState(newState); err != nil {
				log.Printf("Error %s pushing state: %v", err, newState)
				continue
			}

			state = newState

			if err := aircon.writeState(state); err != nil {
				log.Printf("Error persisting state: %s", err)
			}
		}

		if err := aircon.removeConfig(); err != nil {
			log.Printf("Error removing retained config: %s", err)
		}
		if err := aircon.removeState(); err != nil {
			log.Printf("Error remove retained state: %s", err)
		}

		stopDoneChan <- true
	}()

	return func() {
		stopChan <- true
		<-stopDoneChan
	}, nil
}
