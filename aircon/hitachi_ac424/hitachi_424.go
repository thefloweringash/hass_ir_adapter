package hitachi_ac424

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type Mode byte

const (
	ModeHeating Mode = 1
	ModeCooling Mode = 3
)

type DehumidPower byte

const (
	DehumidStrong  DehumidPower = 0
	DehumidCooling DehumidPower = 6
)

func checksum(d []byte) byte {
	var checksum byte
	for _, b := range d {
		checksum += b
	}
	return checksum
}

type State struct {
	On          bool
	Mode        Mode
	Temperature byte
}


type Data struct {

}

func (config FullState) Encode() ([]byte, error) {
	var onOff byte
	if config.On {
		onOff = 1
	}

	var dehumidPower DehumidPower
	if config.Mode == ModeCooling {
		dehumidPower = DehumidCooling
	} else {
		dehumidPower = DehumidStrong
	}

	packet := Packet{
		Header: []byte{0x23, 0xcb, 0x26, 0x01, 0x00},

		OnOff:        onOff,
		Temperature:  16 + config.Temperature,
		Mode:         byte(config.Mode),
		DehumidPower: byte(dehumidPower),
		LeftRight:    3,
		BeepCount:    1,
	}

	packed, err := restruct.Pack(binary.LittleEndian, &packet)
	if err != nil {
		return packed, err
	}

	return append(packed, checksum(packed)), nil
}

