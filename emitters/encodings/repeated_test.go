package encodings

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thefloweringash/hass_ir_adapter/emitters/irkit"
)

func TestRepeated_ToIntervals(t *testing.T) {
	raw := Raw{Intervals: []uint16{1}}

	intervals := Repeat(raw, 35).(irkit.ToIntervals).ToIntervals()

	assert.Equal(t, []uint16{1, 35000, 1}, intervals)
}
