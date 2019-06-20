package intervals

// Adapted from:
// https://github.com/markszabo/IRremoteESP8266/blob/v2.3.3/src/ir_Panasonic.cpp#L19-L32

const (
	PanasonicTick           = 432
	PanasonicHdrMarkTicks   = 8
	PanasonicHdrMark        = PanasonicTick * PanasonicHdrMarkTicks
	PanasonicHdrSpaceTicks  = 4
	PanasonicHdrSpace       = PanasonicTick * PanasonicHdrSpaceTicks
	PanasonicBitMarkTicks   = 1
	PanasonicBitMark        = PanasonicTick * PanasonicBitMarkTicks
	PanasonicOneSpaceTicks  = 3
	PanasonicOneSpace       = PanasonicTick * PanasonicOneSpaceTicks
	PanasonicZeroSpaceTicks = 1
	PanasonicZeroSpace      = PanasonicTick * PanasonicZeroSpaceTicks
)

func EncodePanasonic(payload []byte) []uint16 {
	result := []uint16{PanasonicHdrMark, PanasonicHdrSpace}

	for _, b := range payload {
		for i := uint8(0); i < 8; i++ {
			var space uint16
			if b&(1<<i) != 0 {
				space = PanasonicOneSpace
			} else {
				space = PanasonicZeroSpace
			}
			result = append(result, PanasonicBitMark, space)
		}
	}

	result = append(result, PanasonicBitMark)

	return result
}
