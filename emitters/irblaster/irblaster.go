package irblaster

import (
	"bytes"
	"encoding/binary"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

type MQTTIRBlaster struct {
	client mqtt.Client
	topic  string
}

func NewMQTTIRBlaster(client mqtt.Client, topic string) emitters.Emitter {
	return &MQTTIRBlaster{client, topic}
}

func (self *MQTTIRBlaster) Emit(commands ...emitters.Command) error {
	data := new(bytes.Buffer)

	if err := binary.Write(data, binary.LittleEndian, byte(len(commands))); err != nil {
		return err
	}

	for _, command := range commands {
		encodedParts := []interface{}{
			byte(len(command.Payload) + 1),
			byte(command.Encoding),
			command.Payload,
		}
		for _, encoded := range encodedParts {
			if err := binary.Write(data, binary.LittleEndian, encoded); err != nil {
				return err
			}
		}
	}

	token := self.client.Publish(self.topic, 0, false, data.Bytes())
	token.Wait()
	return token.Error()
}
