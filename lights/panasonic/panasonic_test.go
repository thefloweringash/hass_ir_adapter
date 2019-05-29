package panasonic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Encode(t *testing.T) {
	successfulEncodeTests := []struct {
		Command
		expectedEncode []byte
	}{
		{Command{Channel: Channel1, Action: ActionOn},
			[]byte{0x2c, 0x52, 0x9, 0x2d, 0x24}},
		{Command{Channel: Channel1, Action: ActionOff},
			[]byte{0x2c, 0x52, 0x9, 0x2f, 0x26}},
		{Command{Channel: Channel1, Action: ActionWhite},
			[]byte{0x2c, 0x52, 0x39, 0x90, 0xa9}},
		{Command{Channel: Channel1, Action: ActionWarm},
			[]byte{0x2c, 0x52, 0x39, 0x91, 0xa8}},
		{Command{Channel: Channel1, Action: ActionBright},
			[]byte{0x2c, 0x52, 0x9, 0x2a, 0x23}},
		{Command{Channel: Channel1, Action: ActionDark},
			[]byte{0x2c, 0x52, 0x9, 0x2b, 0x22}},
		{Command{Channel: Channel1, Action: ActionFull},
			[]byte{0x2c, 0x52, 0x9, 0x2c, 0x25}},
		{Command{Channel: Channel1, Action: ActionNight},
			[]byte{0x2c, 0x52, 0x9, 0x2e, 0x27}},
		{Command{Channel: Channel1, Action: ActionSleep},
			[]byte{0x2c, 0x52, 0x39, 0xa1, 0x98}},
		{Command{Channel: Channel2, Action: ActionOn},
			[]byte{0x2c, 0x52, 0x9, 0x35, 0x3c}},
		{Command{Channel: Channel2, Action: ActionOff},
			[]byte{0x2c, 0x52, 0x9, 0x37, 0x3e}},
		{Command{Channel: Channel2, Action: ActionWhite},
			[]byte{0x2c, 0x52, 0x39, 0x94, 0xad}},
		{Command{Channel: Channel2, Action: ActionWarm},
			[]byte{0x2c, 0x52, 0x39, 0x95, 0xac}},
		{Command{Channel: Channel2, Action: ActionBright},
			[]byte{0x2c, 0x52, 0x9, 0x32, 0x3b}},
		{Command{Channel: Channel2, Action: ActionDark},
			[]byte{0x2c, 0x52, 0x9, 0x33, 0x3a}},
		{Command{Channel: Channel2, Action: ActionFull},
			[]byte{0x2c, 0x52, 0x9, 0x34, 0x3d}},
		{Command{Channel: Channel2, Action: ActionNight},
			[]byte{0x2c, 0x52, 0x9, 0x36, 0x3f}},
		{Command{Channel: Channel2, Action: ActionSleep},
			[]byte{0x2c, 0x52, 0x39, 0xaa, 0x93}},
		{Command{Channel: Channel3, Action: ActionOn},
			[]byte{0x2c, 0x52, 0x9, 0x3d, 0x34}},
		{Command{Channel: Channel3, Action: ActionOff},
			[]byte{0x2c, 0x52, 0x9, 0x3f, 0x36}},
		{Command{Channel: Channel3, Action: ActionWhite},
			[]byte{0x2c, 0x52, 0x39, 0x98, 0xa1}},
		{Command{Channel: Channel3, Action: ActionWarm},
			[]byte{0x2c, 0x52, 0x39, 0x99, 0xa0}},
		{Command{Channel: Channel3, Action: ActionBright},
			[]byte{0x2c, 0x52, 0x9, 0x3a, 0x33}},
		{Command{Channel: Channel3, Action: ActionDark},
			[]byte{0x2c, 0x52, 0x9, 0x3b, 0x32}},
		{Command{Channel: Channel3, Action: ActionFull},
			[]byte{0x2c, 0x52, 0x9, 0x3c, 0x35}},
		{Command{Channel: Channel3, Action: ActionNight},
			[]byte{0x2c, 0x52, 0x9, 0x3e, 0x37}},
		{Command{Channel: Channel3, Action: ActionSleep},
			[]byte{0x2c, 0x52, 0x39, 0xb3, 0x8a}},
	}

	for _, test := range successfulEncodeTests {
		actualEncode, err := test.Command.Encode()
		assert.Nil(t, err)
		assert.Equal(t, test.expectedEncode, actualEncode,
			fmt.Sprintf("Encoding Command: %#v", test.Command))

	}
}
