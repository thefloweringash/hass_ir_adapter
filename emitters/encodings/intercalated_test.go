package encodings

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/emitters/irblaster"
	"github.com/thefloweringash/hass_ir_adapter/emitters/irkit"
)

func TestIntercalated_ToIntervals(t *testing.T) {
	raws := []emitters.Command{
		Raw{Intervals: []uint16{1}},
		Raw{Intervals: []uint16{2}},
		Raw{Intervals: []uint16{3}},
	}

	intervals := Intercalate(raws, 35).(irkit.ToIntervals).ToIntervals()

	assert.Equal(t, []uint16{1, 35000, 2, 35000, 3}, intervals)
}

func TestIntercalated_ToIRBlasterEncoding(t *testing.T) {
	raws := []emitters.Command{
		Raw{Intervals: []uint16{1}},
		Raw{Intervals: []uint16{2}},
		Raw{Intervals: []uint16{3}},
	}

	commands := Intercalate(raws, 35).(irblaster.ToIRBlasterEncoding).ToIRBlasterEncoding()

	assert.Equal(t, []irblaster.Command{
		{Encoding: EncodingRaw, Payload: []interface{}{uint32(38), []uint16{1}}},
		{Encoding: EncodingDelay, Payload: []interface{}{uint16(35)}},
		{Encoding: EncodingRaw, Payload: []interface{}{uint32(38), []uint16{2}}},
		{Encoding: EncodingDelay, Payload: []interface{}{uint16(35)}},
		{Encoding: EncodingRaw, Payload: []interface{}{uint32(38), []uint16{3}}},
	}, commands)
}
