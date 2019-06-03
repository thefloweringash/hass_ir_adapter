package daiko

import (
	"fmt"
)

type commandCode int

const (
	commandOff    commandCode = 0
	commandToggle commandCode = 5
	commandWhite  commandCode = 40
	commandFull   commandCode = 41
	commandWarm   commandCode = 42
)

type Channel uint8

const (
	Channel1 Channel = 0
	Channel2 Channel = 1
)

const (
	BrightnessMin = 1
	BrightnessMax = 10
	WarmthMin     = 1
	WarmthMax     = 11
)

var brightnessLimits = map[uint8]uint8{
	1:  7,
	2:  7,
	3:  8,
	4:  8,
	5:  9,
	6:  10,
	7:  9,
	8:  8,
	9:  8,
	10: 7,
	11: 7,
}

type limits struct {
	min, max uint8
}

var warmthLimits = map[uint8]limits{
	7:  {1, 11},
	8:  {3, 9},
	9:  {5, 7},
	10: {6, 6},
}

func ClampBrightness(warmth uint8, brightness uint8) uint8 {
	maxBrightness := brightnessLimits[warmth]
	if brightness <= maxBrightness {
		return brightness
	} else {
		return maxBrightness
	}
}

func ClampWarmth(warmth uint8, brightness uint8) uint8 {
	limit, haveLimit := warmthLimits[brightness]
	switch {
	case !haveLimit:
		return warmth
	case warmth < limit.min:
		return limit.min
	case warmth > limit.max:
		return limit.max
	default:
		return warmth
	}
}

var header = []byte{0x85, 0xfb}

func pack(dest []byte, channel Channel, payload uint8) {
	dest[0] = payload
	dest[0] |= byte(channel) << 7
	dest[1] = ^dest[0]
}

func shortPacket(channel Channel, b1 uint8) []byte {
	result := make([]byte, 4)
	copy(result, header)
	pack(result[2:4], channel, b1)
	return result
}

func longPacket(channel Channel, b1 uint8, b2 uint8) []byte {
	result := make([]byte, 6)
	copy(result, header)
	pack(result[2:4], channel, b1)
	pack(result[4:6], channel, b2)
	return result
}

func Off(channel Channel) []byte {
	return shortPacket(channel, uint8(commandOff))
}

func Toggle(channel Channel) []byte {
	return shortPacket(channel, uint8(commandToggle))
}

func White(channel Channel) []byte {
	return shortPacket(channel, uint8(commandWhite))
}

func Full(channel Channel) []byte {
	return shortPacket(channel, uint8(commandFull))
}

func Warm(channel Channel) []byte {
	return shortPacket(channel, uint8(commandWarm))
}

func NightLight(channel Channel, intensity uint8) ([]byte, error) {
	// TODO: bounds
	var val uint8
	if intensity <= 7 {
		val = intensity + 5
	} else {
		val = intensity + 8
	}
	return shortPacket(channel, val), nil
}

func On(channel Channel, warmth uint8, brightness uint8) ([]byte, error) {
	if !(warmth >= WarmthMin && warmth <= WarmthMax) {
		return nil, fmt.Errorf("warmth %d outside range %d...%d",
			warmth, WarmthMin, WarmthMax)
	}
	if !(brightness >= BrightnessMin && brightness <= BrightnessMax) {
		return nil, fmt.Errorf("brightness %d outside range %d...%d",
			brightness, BrightnessMin, BrightnessMax)
	}

	maxBrightness := brightnessLimits[warmth]

	if brightness > maxBrightness {
		err := fmt.Errorf("brightness %d is greater than maximum %d for warmth %d",
			brightness, maxBrightness, warmth)
		return nil, err
	}

	return longPacket(channel, warmth+28, brightness+18), nil
}
