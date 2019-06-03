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

var EncodingPanasonic = GenericEncoding{
	HeaderMark:  PanasonicHdrMark,
	HeaderSpace: PanasonicHdrSpace,
	BitMark:     PanasonicBitMark,
	OneSpace:    PanasonicOneSpace,
	ZeroSpace:   PanasonicZeroSpace,
}

func EncodePanasonic(payload []byte) []uint16 {
	return EncodeGeneric(EncodingPanasonic, payload)
}
