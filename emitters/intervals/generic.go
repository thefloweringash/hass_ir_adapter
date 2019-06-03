package intervals

type GenericEncoding struct {
	HeaderMark, HeaderSpace      uint16
	BitMark, OneSpace, ZeroSpace uint16
}

func EncodeGeneric(encoding GenericEncoding, payload []byte) []uint16 {
	result := []uint16{encoding.HeaderMark, encoding.HeaderSpace}

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
