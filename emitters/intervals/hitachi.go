package intervals

// Adapted from:
// https://github.com/crankyoldgit/IRremoteESP8266/blob/0683f123a0c4fdd4f71e74d54a972b109c028335/src/ir_Hitachi.cpp#L32-L39

const (
	HitachiAc424LdrMark  = 29784
	HitachiAc424LdrSpace = 49290

	HitachiAc424HdrMark   = 3416
	HitachiAc424HdrSpace  = 1604
	HitachiAc424BitMark   = 463
	HitachiAc424OneSpace  = 1208
	HitachiAc424ZeroSpace = 372
)

var encodingHitachiAc424Body = GenericEncoding{
	HeaderMark:  HitachiAc424HdrMark,
	HeaderSpace: HitachiAc424HdrSpace,
	BitMark:     HitachiAc424BitMark,
	OneSpace:    HitachiAc424OneSpace,
	ZeroSpace:   HitachiAc424ZeroSpace,
}

type encodingHitachiAc424 struct{}

var hitatchAc424Leader = []uint16{HitachiAc424LdrMark, HitachiAc424LdrSpace}

func (encoding encodingHitachiAc424) Encode(payload []byte) []uint16 {
	return append(hitatchAc424Leader, encodingHitachiAc424Body.Encode(payload)...)
}

var EncodingHitachiAc424 = encodingHitachiAc424{}
