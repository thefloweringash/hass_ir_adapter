package device

import (
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Mqtt struct {
	Class  string
	Id     string
	Client mqtt.Client
	Logger *log.Logger
}

func (d Mqtt) Prefix() string {
	return fmt.Sprintf("homeassistant/%s/%s", d.Class, d.Id)
}

func (d Mqtt) TopicName(topic string) string {
	return fmt.Sprintf("%s/%s", d.Prefix(), topic)
}

func (d Mqtt) Subscribe(channel chan Update, state State) error {
	for _, binding := range state.Bindings() {
		b := binding

		handler := func(c mqtt.Client, m mqtt.Message) {
			channel <- Update{binding: b, value: string(m.Payload())}
		}

		topic := d.TopicName(binding.RelativeTopic())
		token := d.Client.Subscribe(topic, 0, handler)
		token.Wait()
		if err := token.Error(); err != nil {
			return err
		}
		d.Logger.Printf("subscribed to %s", topic)
	}

	return nil
}

const (
	RelativeConfigTopic = "config"
	RelativeStateTopic  = "state"
)

func (d Mqtt) publishJson(relativeTopic string, value interface{}) error {
	value, err := json.Marshal(value)
	if err != nil {
		return err
	}

	token := d.Client.Publish(d.TopicName(relativeTopic), 0, true, value)
	token.Wait()
	return token.Error()
}

func (d Mqtt) PublishConfig(name string, state State, config map[string]interface{}) error {
	configMap := config

	for k, v := range GenerateStateConfig(name, d.Prefix(), state) {
		configMap[k] = v
	}

	return d.publishJson(RelativeConfigTopic, configMap)
}

func (d Mqtt) PublishState(state State) error {
	return d.publishJson(RelativeStateTopic, state)
}

func (d Mqtt) RemoveConfig() error {
	token := d.Client.Publish(d.TopicName(RelativeConfigTopic), 0, true, []byte{})
	token.Wait()
	return token.Error()
}

func (d Mqtt) RemoveState() error {
	token := d.Client.Publish(d.TopicName(RelativeStateTopic), 0, true, []byte{})
	token.Wait()
	return token.Error()
}
