package encodings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeated_ToIntervals(t *testing.T) {
	raw := Raw{Intervals: []uint16{1}}

	intervals := Repeat(raw, 35).(ToIntervals).ToIntervals()

	assert.Equal(t, []uint16{1, 35000, 1}, intervals)
}
