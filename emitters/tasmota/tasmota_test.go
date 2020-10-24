package tasmota

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_compactEncoding(t *testing.T) {
	encodingExamples := []struct {
		input   []uint16
		encoded string
	}{
		// Minimal example
		{[]uint16{432, 0, 432, 0},
			"+432-0Ab"},

		// Example from the website
		{[]uint16{8570, 4240, 550, 1580, 550, 510, 565, 1565, 565, 505, 565, 505},
			"+8570-4240+550-1580C-510+565-1565F-505Fh"},
	}

	for _, example := range encodingExamples {
		assert.Equal(t, example.encoded, compactEncoding(example.input))
	}
}
