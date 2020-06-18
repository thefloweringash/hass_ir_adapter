package encodings

import (
	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/emitters/intervals"
	"github.com/thefloweringash/hass_ir_adapter/emitters/irblaster"
)

type daikinHeader struct{}

var daikinHeaderIntervals = intervals.EncodingPanasonic.EncodeBits(0, 5)

func (hdr daikinHeader) ToIRBlasterEncoding() []irblaster.Command {
	return Raw{Intervals: hdr.ToIntervals()}.ToIRBlasterEncoding()
}

func (daikinHeader) ToIntervals() []uint16 {
	result := make([]uint16, len(daikinHeaderIntervals))
	copy(result, daikinHeaderIntervals)
	return result
}

type Daikin struct {
	P1, P2 []uint8
}

func (daikin Daikin) build() emitters.Command {
	return Intercalate([]emitters.Command{
		daikinHeader{},
		Panasonic{daikin.P1},
		Panasonic{daikin.P2},
	}, 29)

}

func (daikin Daikin) ToIRBlasterEncoding() []irblaster.Command {
	return daikin.build().(irblaster.ToIRBlasterEncoding).ToIRBlasterEncoding()
}

func (daikin Daikin) ToIntervals() []uint16 {
	return daikin.build().(ToIntervals).ToIntervals()
}
