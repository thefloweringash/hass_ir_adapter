package emitters

import (
	"github.com/eclipse/paho.mqtt.golang"
)

type Emitter interface {
	Emit(commands ...Command) error
}

const (
	Panasonic Encoding = 241
)

type Encoding int

type Command struct {
	Encoding Encoding
	Payload  []byte
}

type Factory interface {
	New(c mqtt.Client) (Emitter, error)
	Id() string
}
