package intervals

// Adapted from:
// https://github.com/markszabo/IRremoteESP8266/blob/v2.3.3/src/ir_NEC.cpp#L22-L32

const (
	NecTick           = 560
	NecHdrMarkTicks   = 16
	NecHdrMark        = NecTick * NecHdrMarkTicks
	NecHdrSpaceTicks  = 8
	NecHdrSpace       = NecTick * NecHdrSpaceTicks
	NecBitMarkTicks   = 1
	NecBitMark        = NecTick * NecBitMarkTicks
	NecOneSpaceTicks  = 3
	NecOneSpace       = NecTick * NecOneSpaceTicks
	NecZeroSpaceTicks = 1
	NecZeroSpace      = NecTick * NecZeroSpaceTicks
)

var EncodingNec = GenericEncoding{
	HeaderMark:  NecHdrMark,
	HeaderSpace: NecHdrSpace,
	BitMark:     NecBitMark,
	OneSpace:    NecOneSpace,
	ZeroSpace:   NecZeroSpace,
}
