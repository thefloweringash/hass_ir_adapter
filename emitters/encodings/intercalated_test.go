package encodings

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thefloweringash/hass_ir_adapter/emitters"
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
