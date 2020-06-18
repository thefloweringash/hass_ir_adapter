package tasmota

import (
	"errors"
	"log"
	"math"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/emitters/encodings"
)

type TasmotaEmitter struct {
	client mqtt.Client
	logger *log.Logger
	topic  string
}

func (self *TasmotaEmitter) Lock() {

}

func (self *TasmotaEmitter) Unlock() {

}

func NewTasmotaEmitter(client mqtt.Client, logger *log.Logger, topic string) emitters.Emitter {
	return &TasmotaEmitter{client, logger, topic}
}

func (self *TasmotaEmitter) Emit(command emitters.Command) error {
	intervalCommand, ok := command.(encodings.ToIntervals)
	if !ok {
		return errors.New("command not convertible to intervals")
	}

	intervals := intervalCommand.ToIntervals()
	tasmotaIntervals := []string{"0"}

	self.logger.Printf("raw intervals: %v\n", intervals)

	for i := range intervals {
		x := intervals[i]
		for x > math.MaxUint16 {
			tasmotaIntervals = append(tasmotaIntervals, strconv.Itoa(math.MaxUint16), "0")
			x -= math.MaxUint16
		}
		tasmotaIntervals = append(tasmotaIntervals, strconv.Itoa(int(x)))
	}

	data := strings.Join(tasmotaIntervals, ",")

	self.logger.Printf("command: %v\n", data)

	token := self.client.Publish(self.topic+"/IRsend", 0, false, data)
	token.Wait()
	return token.Error()
}
