package daikin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState_Encode(t *testing.T) {
	successfulEncodeTests := []struct {
		state                  State
		expectedp1, expectedp2 []byte
	}{
		{
			state: State{
				On:          false,
				Temperature: 18,
			},
			expectedp1: []byte{
				0x11, 0xda, 0x27, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x14,
			},
			expectedp2: []byte{
				0x11, 0xda, 0x27, 0x0, 0x0, 0x8, 0x24, 0x0, 0x3f, 0x0, 0x0, 0x6, 0x60, 0x0, 0x0, 0xc3, 0x0, 0x0, 0xa6,
			},
		},
	}

	for _, test := range successfulEncodeTests {
		p1, p2, err := test.state.Encode()
		assert.Nil(t, err)
		assert.Equal(t, test.expectedp1, p1)
		assert.Equal(t, test.expectedp2, p2)
		assert.Equal(t, 20, len(p1))
		assert.Equal(t, 19, len(p2))
	}
}
