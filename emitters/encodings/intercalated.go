package encodings

import (
	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/emitters/irblaster"
	"github.com/thefloweringash/hass_ir_adapter/emitters/irkit"
)

const (
	EncodingDelay = 239
)

type intercalated struct {
	inner []emitters.Command
	gap   uint16
}

func (r intercalated) ToIntervals() []uint16 {
	var intervals []uint16

	for _, i := range r.inner {
		inner, ok := i.(irkit.ToIntervals)
		if !ok {
			panic("Can't intercalate via intervals unless the repeated encoding is intervals-able!")
		}

		intervals = append(intervals, inner.ToIntervals()...)
		intervals = append(intervals, r.gap*1000)
	}

	return intervals[:len(intervals)-1]
}

func (r intercalated) ToIRBlasterEncoding() []irblaster.Command {
	var commands []irblaster.Command

	delay := irblaster.Command{
		Encoding: EncodingDelay,
		Payload:  []interface{}{r.gap},
	}

	for _, i := range r.inner {
		inner, ok := i.(irblaster.ToIRBlasterEncoding)
		if !ok {
			panic("Can't intercalate via irblaster unless the repeated encoding is irblaster-able!")
		}

		commands = append(commands, inner.ToIRBlasterEncoding()...)
		commands = append(commands, delay)
	}
	return commands[:len(commands)-1]
}

func Intercalate(commands []emitters.Command, gap uint16) emitters.Command {
	return intercalated{inner: commands, gap: gap}
}
