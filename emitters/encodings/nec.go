package encodings

import (
	"github.com/thefloweringash/hass_ir_adapter/emitters/intervals"
	"github.com/thefloweringash/hass_ir_adapter/emitters/irblaster"
)

type NEC struct {
	Payload []byte
}

func (command NEC) ToIntervals() []uint16 {
	return intervals.EncodeNec(command.Payload)
}

func (command NEC) ToIRBlasterEncoding() []irblaster.Command {
	return []irblaster.Command{{
		Encoding: EncodingRaw,
		Payload: []interface{}{
			uint32(38),
			command.ToIntervals(),
		},
	}}
}
