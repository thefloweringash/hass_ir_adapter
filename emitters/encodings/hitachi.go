package encodings

import (
	"github.com/thefloweringash/hass_ir_adapter/emitters/intervals"
	"github.com/thefloweringash/hass_ir_adapter/emitters/irblaster"
)

type HitachiAc424 struct {
	Payload []byte
}

func (command HitachiAc424) ToIntervals() []uint16 {
	return intervals.EncodingHitachiAc424.Encode(command.Payload)
}

func (command HitachiAc424) ToIRBlasterEncoding() []irblaster.Command {
	return intervalsToIRBlasterCommand(command.ToIntervals())
}
