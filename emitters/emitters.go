package emitters

import (
	"log"

	"github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/emitters/intervals"
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

func (cmd Command) Intervals() []uint16 {
	switch cmd.Encoding {
	case Panasonic:
		return intervals.EncodePanasonic(cmd.Payload)
	}
	panic("Encoding unknown command type")
}

type Factory interface {
	New(c mqtt.Client, logger *log.Logger) (Emitter, error)
	Id() string
}
