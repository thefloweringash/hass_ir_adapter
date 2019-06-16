package panasonic

import (
	"github.com/thefloweringash/hass_ir_adapter/emitters"
	"github.com/thefloweringash/hass_ir_adapter/lights"
)

type Device struct {
	emitter emitters.Emitter
	channel Channel
}

func (device *Device) PushState(state lights.State) error {
	command := Command{Channel: device.channel}
	switch {
	case !state.On:
		command.Action = ActionOff
	case state.Brightness < 26:
		command.Action = ActionNight
	case state.Brightness >= 230:
		command.Action = ActionFull
	default:
		command.Action = ActionOn
	}

	encoded, err := command.Encode()
	if err != nil {
		return err
	}

	return device.emitter.Emit(emitters.Command{
		Encoding: emitters.Panasonic,
		Payload:  encoded,
	})
}
