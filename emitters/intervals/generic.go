package intervals

type GenericEncoding struct {
	HeaderMark, HeaderSpace      uint16
	BitMark, OneSpace, ZeroSpace uint16
}

func (encoding GenericEncoding) EncodeHeader() []uint16 {
	return []uint16{encoding.HeaderMark, encoding.HeaderSpace}
}

func (encoding GenericEncoding) EncodeBody(payload []byte) []uint16 {
	var result []uint16

	for _, b := range payload {
		for i := uint8(0); i < 8; i++ {
			var space uint16
			if b&(1<<i) != 0 {
				space = encoding.OneSpace
			} else {
				space = encoding.ZeroSpace
			}
			result = append(result, encoding.BitMark, space)
		}
	}

	result = append(result, encoding.BitMark)

	return result
}

func (encoding GenericEncoding) EncodeBits(payload byte, bits uint8) []uint16 {
	var result []uint16

	for i := uint8(0); i < bits; i++ {
		var space uint16
		if payload&(1<<i) != 0 {
			space = encoding.OneSpace
		} else {
			space = encoding.ZeroSpace
		}
		result = append(result, encoding.BitMark, space)
	}

	result = append(result, encoding.BitMark)

	return result
}

func (encoding GenericEncoding) Encode(payload []byte) []uint16 {
	return append(encoding.EncodeHeader(), encoding.EncodeBody(payload)...)
}
