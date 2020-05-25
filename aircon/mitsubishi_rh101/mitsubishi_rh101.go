package mitsubishi_rh101

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

type FullState struct {
	On          bool
	Mode        Mode
	Temperature byte
}

type Packet struct {
	Header []byte `struct:"[5]byte"`

	Padding1 byte `struct:"uint8:2"`
	OnOff    byte `struct:"uint8:1"`
	Padding2 byte `struct:"uint8:5"`

	Padding3 byte `struct:"uint8:2"`
	Mode     byte `struct:"uint8:3"`
	Padding4 byte `struct:"uint8:3"`

	Padding5    byte `struct:"uint8:4"`
	Temperature byte `struct:"uint8:4"`

	LeftRight    byte `struct:"uint8:4"`
	DehumidPower byte `struct:"uint8:4"`

	BeepCount     byte `struct:"uint8:2"`
	WindDirection byte `struct:"uint8:3"`
	WindSpeed     byte `struct:"uint8:3"`

	Clock byte `struct:"uint8"`

	EndTime   byte `struct:"uint8"`
	StartTime byte `struct:"uint8"`

	Padding6 byte `struct:"uint8:5"`
	ProgMode byte `struct:"uint8:1"`

	Padding7 byte `struct:"uint8"`

	Padding8 byte `struct:"uint8:3"`
	Powerful byte `struct:"uint8:1"`
	Padding9 byte `struct:"uint8:4"`

	Padding10 byte `struct:"uint8"`
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
