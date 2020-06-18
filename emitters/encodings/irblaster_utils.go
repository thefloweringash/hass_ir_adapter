package encodings

import (
	"github.com/thefloweringash/hass_ir_adapter/emitters/irblaster"
)

func intervalsToIRBlasterCommand(intervals []uint16) []irblaster.Command {
	return []irblaster.Command{{
		Encoding: EncodingRaw,
		Payload: []interface{}{
			uint32(38),
			intervals,
		},
	}}
}
