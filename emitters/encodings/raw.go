package encodings

import (
	"github.com/thefloweringash/hass_ir_adapter/emitters/irblaster"
)

type Raw struct {
	Intervals []uint16
}

func (raw Raw) ToIntervals() []uint16 {
	return raw.Intervals
}

const EncodingRaw = 240

func (raw Raw) ToIRBlasterEncoding() []irblaster.Command {
	return []irblaster.Command{{
		Encoding: EncodingRaw,
		Payload: []interface{}{
			uint32(38),
			raw.Intervals,
		}},
	}
}
