package tasmota

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/emitters/encodings"
)

type TasmotaEmitter struct {
	client  mqtt.Client
	logger  *log.Logger
	topic   string
	compact bool
}

func (self *TasmotaEmitter) Lock() {

}

func (self *TasmotaEmitter) Unlock() {

}

func NewTasmotaEmitter(client mqtt.Client, logger *log.Logger, topic string, compact bool) emitters.Emitter {
	return &TasmotaEmitter{client, logger, topic, compact}
}

func compactEncoding(intervals []uint16) string {
	var result strings.Builder
	names := map[uint16]string{}

	nextName := 'A'
	emitHigh := true

	for _, interval := range intervals {
		if name, hasName := names[interval]; hasName {
			if emitHigh {
				fmt.Fprint(&result, name)
			} else {
				fmt.Fprint(&result, strings.ToLower(name))
			}
		} else {
			if emitHigh {
				fmt.Fprintf(&result, "+%d", interval)
			} else {
				fmt.Fprintf(&result, "-%d", interval)
			}

			names[interval] = string(nextName)
			nextName = nextName + 1
		}

		emitHigh = !emitHigh
	}

	return result.String()
}

func simpleEncoding(intervals []uint16) string {
	var stringIntervals []string
	for _, interval := range intervals {
		stringIntervals = append(stringIntervals, strconv.Itoa(int(interval)))
	}
	return strings.Join(stringIntervals, ",")
}

func (self *TasmotaEmitter) Emit(command emitters.Command) error {
	intervalCommand, ok := command.(encodings.ToIntervals)
	if !ok {
		return errors.New("command not convertible to intervals")
	}

	intervals := intervalCommand.ToIntervals()
	tasmotaIntervals := []uint16{}

	self.logger.Printf("raw intervals: %v\n", intervals)

	for i := range intervals {
		x := intervals[i]
		for x > math.MaxUint16 {
			tasmotaIntervals = append(tasmotaIntervals, math.MaxUint16, 0)
			x -= math.MaxUint16
		}
		tasmotaIntervals = append(tasmotaIntervals, x)
	}

	var encodedIntervals string
	if self.compact {
		encodedIntervals = compactEncoding(intervals)
	} else {
		encodedIntervals = simpleEncoding(intervals)
	}

	data := "0," + encodedIntervals

	self.logger.Printf("command: %v\n", data)

	token := self.client.Publish(self.topic+"/IRsend", 0, false, data)
	token.Wait()
	return token.Error()
}
