package daikin

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type Mode uint8

const (
	ModeAuto        Mode = 0
	ModeDry         Mode = 2
	ModeCooling     Mode = 3
	ModeHeating     Mode = 4
	ModeSendoffWind Mode = 6
)

func checksum(buf []byte) uint8 {
	var c uint8
	for _, x := range buf {
		c += x
	}
	return c
}

func encodeTimers(onTimer, offTimer uint16) []byte {
	var onTimerLow, onTimerHigh byte
	var offTimerLow, offTimerHigh byte

	if onTimer != 0 {
		onTimerLow = uint8(onTimer & 0xff)
		onTimerHigh = uint8((onTimer >> 8) & 0x7)
	} else {
		onTimerHigh = 3 << 1
	}

	if offTimer != 0 {
		offTimerLow = uint8(offTimer & 0xf)
		offTimerHigh = uint8((offTimer >> 4) & 0x7f)
	} else {
		offTimerHigh = 3 << 5
	}

	buf := make([]byte, 3)
	buf[0] = onTimerLow
	buf[1] = offTimerLow<<4 | onTimerHigh
	buf[2] = offTimerHigh
	return buf
}

var (
	p1header = []byte{17, 218, 39, 0, 2}
	p2header = []byte{17, 218, 39, 0, 0}
)

type p1 struct {
	Header    []byte `struct:"[5]byte"`
	Padding   []byte `struct:"[7]byte"`
	FanYonder bool   `struct:"uint8:1"`
	Rest      []byte `struct:"[6]byte"`
}

type p2 struct {
	Header []byte `struct:"[5]byte"`

	Padding1 uint8 `struct:"uint8:1"`
	Mode     Mode  `struct:"uint8:3"`
	Padding2 uint8 `struct:"uint8:1"`
	OffTimer bool  `struct:"uint8:1"`
	OnTimer  bool  `struct:"uint8:1"`
	Power    bool  `struct:"uint8:1"`

	Temperature uint8

	Padding3 uint8

	FanSpeed      uint8 `struct:"uint8:4"`
	VaneDirection uint8 `struct:"uint8:4"`

	LeftRight uint8 `struct:"uint8:4"`
	Padding4  uint8 `struct:"uint8:4"`

	Timers []byte `struct:"[3]byte"`

	Padding5     uint8 `struct:"uint8:3"`
	Silent       bool  `struct:"uint8:1"`
	Padding6     uint8 `struct:"uint8:3"`
	PowerfulMode bool  `struct:"uint8:1"`

	IntelligentOn bool  `struct:"uint8:1"`
	Padding7      uint8 `struct:"uint8:7"`

	Padding8 uint8
	Padding9 uint16
}

type State struct {
	On          bool
	Mode        Mode
	Temperature float32
}

func (state State) Encode() ([]byte, []byte, error) {
	p1 := p1{
		Header:    p1header,
		FanYonder: false,
	}

	p1bytes, err := restruct.Pack(binary.LittleEndian, &p1)
	if err != nil {
		return nil, nil, err
	}

	p1bytes = append(p1bytes, checksum(p1bytes))

	p2 := p2{
		Header:        p2header,
		Mode:          state.Mode,
		Power:         state.On,
		Padding2:      1,
		Temperature:   uint8(state.Temperature * 2),
		FanSpeed:      3,
		VaneDirection: 15,
		Timers:        encodeTimers(0, 0),
		Padding8:      195,
	}

	p2bytes, err := restruct.Pack(binary.LittleEndian, &p2)
	if err != nil {
		return nil, nil, err
	}

	p2bytes = append(p2bytes, checksum(p2bytes))

	return p1bytes, p2bytes, nil
}
