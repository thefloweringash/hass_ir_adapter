package irblaster

import (
	"bytes"
	"encoding/binary"
	"errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

type ToIRBlasterEncoding interface {
	ToIRBlasterEncoding() []Command
}

type Command struct {
	Encoding uint8
	Payload  []interface{}
}

type MQTTIRBlaster struct {
	client mqtt.Client
	topic  string
}

// No meaningful locking. The emitter will queue in mqtt, and we don't
// have feedback for when a queued even has been processed.
func (self *MQTTIRBlaster) Lock() {
}

func (self *MQTTIRBlaster) Unlock() {
}

func NewMQTTIRBlaster(client mqtt.Client, topic string) emitters.Emitter {
	return &MQTTIRBlaster{client, topic}
}

func (self *MQTTIRBlaster) Emit(command emitters.Command) error {
	data := new(bytes.Buffer)

	encodable, ok := command.(ToIRBlasterEncoding)
	if !ok {
		return errors.New("unencodable command")
	}

	commands := encodable.ToIRBlasterEncoding()

	if err := binary.Write(data, binary.LittleEndian, byte(len(commands))); err != nil {
		return err
	}

	for _, command := range commands {
		commandBuffer := new(bytes.Buffer)

		for _, encoded := range command.Payload {
			if err := binary.Write(commandBuffer, binary.LittleEndian, encoded); err != nil {
				return err
			}
		}

		commandBytes := commandBuffer.Bytes()

		encodedParts := []interface{}{
			byte(len(commandBytes) + 1),
			command.Encoding,
			commandBytes,
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
