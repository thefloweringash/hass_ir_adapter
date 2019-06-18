package encodings

import (
	"github.com/thefloweringash/hass_ir_adapter/emitters"
)

func Repeat(command emitters.Command, gap uint16) emitters.Command {
	inners := []emitters.Command{
		command,
		command,
	}
	return Intercalate(inners, gap)
}
