package encodings

import (
	"github.com/thefloweringash/hass_ir_adapter/emitters/intervals"
	"github.com/thefloweringash/hass_ir_adapter/emitters/irblaster"
)

type Panasonic struct {
	Payload []byte
}

func (command Panasonic) ToIntervals() []uint16 {
	return intervals.EncodingPanasonic.Encode(command.Payload)
}

const EncodingPanasonic = 241

func (command Panasonic) ToIRBlasterEncoding() []irblaster.Command {
	return []irblaster.Command{{
		Encoding: EncodingPanasonic,
		Payload:  []interface{}{command.Payload},
	}}
}
