package encodings

import (
	"github.com/thefloweringash/hass_ir_adapter/emitters/irblaster"
	"github.com/thefloweringash/hass_ir_adapter/emitters/irkit"
)

type Repeated struct {
	Inner interface{}
	Gap   uint16
}

func (r Repeated) ToIntervals() []uint16 {
	inner, ok := r.Inner.(irkit.ToIntervals)
	if !ok {
		panic("Can't repeat via intervals unless the repeated encoding is intervals-able!")
	}

	baseIntervals := inner.ToIntervals()
	intervals := append(baseIntervals, r.Gap*1000)
	intervals = append(intervals, baseIntervals...)
	return intervals
}

const (
	EncodingDelay = 239
)

func (r Repeated) ToIRBlasterEncoding() []irblaster.Command {
	inner, ok := r.Inner.(irblaster.ToIRBlasterEncoding)
	if !ok {
		panic("Can't repeat via irblaster unless the repeated encoding is irblaster-able!")
	}

	baseCommands := inner.ToIRBlasterEncoding()
	delay := irblaster.Command{
		Encoding: EncodingDelay,
		Payload:  []interface{}{r.Gap * 1000},
	}
	commands := append(baseCommands, delay)
	commands = append(commands, baseCommands...)
	return commands
}
