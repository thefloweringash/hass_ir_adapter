package emitters

import (
	"log"

	"github.com/eclipse/paho.mqtt.golang"
)

type Emitter interface {
	Emit(commands Command) error
	Lock()
	Unlock()
}

type Command interface {
}

type Delay struct {
	Millis uint16
}

type Raw struct {
	Intervals []uint16
}

type Factory interface {
	New(c mqtt.Client, logger *log.Logger) (Emitter, error)
	Id() string
}
