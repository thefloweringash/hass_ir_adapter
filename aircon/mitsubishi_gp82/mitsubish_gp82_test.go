package mitsubishi_gp82

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thefloweringash/hass_ir_adapter/aircon"
)

func TestMitsubishiAircon_Encode(t *testing.T) {
	successfulEncodeTests := []struct {
		state          aircon.State
		expectedEncode []byte
	}{
		{
			aircon.State{
				Mode:        "cool",
				FanMode:     "auto",
				Temperature: 18,
			},
			[]byte{0x23, 0xcb, 0x26, 0x1, 0x0, 0x24, 0x3, 0xd, 0x0, 0x0, 0x0, 0x0, 0x0, 0x49},
		},
		{
			aircon.State{
				Mode:        "heat",
				FanMode:     "quiet",
				Temperature: 31,
			},
			[]byte{0x23, 0xcb, 0x26, 0x1, 0x0, 0x24, 0x1, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x3c},
		},
	}

	for _, test := range successfulEncodeTests {
		actualEncoded, err := Encode(test.state)
		assert.Nil(t, err)
		assert.Equal(t, test.expectedEncode, actualEncoded)
	}
}

func TestFullState_Encode(t *testing.T) {
	successfulEncodeTests := []struct {
		state          FullState
		expectedEncode []byte
	}{
		{
			FullState{
				On:          true,
				Mode:        ModeCooling,
				Temperature: 18,
			},
			[]byte{0x23, 0xcb, 0x26, 0x1, 0x0, 0x24, 0x3, 0xd, 0x0, 0x0, 0x0, 0x0, 0x0, 0x49},
		},
		{
			FullState{
				On:          true,
				Mode:        ModeHeating,
				WindSpeed:   WindQuiet,
				Temperature: 31,
			},
			[]byte{0x23, 0xcb, 0x26, 0x1, 0x0, 0x24, 0x1, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x3c},
		},
	}

	for _, test := range successfulEncodeTests {
		actualEncoded, err := test.state.Encode()
		assert.Nil(t, err)
		assert.Equal(t, test.expectedEncode, actualEncoded)
	}
}
