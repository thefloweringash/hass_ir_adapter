package panasonic

import (
	"encoding/binary"
	"fmt"

	"gopkg.in/restruct.v1"
)

type Channel int

const (
	Channel1 Channel = 1
	Channel2 Channel = 2
	Channel3 Channel = 3
)

type Action int

const (
	ActionBright = 2
	ActionDark   = 3
	ActionFull   = 4
	ActionOn     = 5
	ActionNight  = 6
	ActionOff    = 7

	ActionWhite = iota + 0x80
	ActionWarm
	ActionSleep
)

type Series int

const (
	Series09 = 0x09
	Series39 = 0x39
)

func (c Action) series() Series {
	switch c {
	case ActionFull, ActionBright, ActionNight, ActionOn, ActionDark, ActionOff:
		return Series09
	default:
		return Series39
	}
}

type Command struct {
	Channel
	Action
}

type Packet09 struct {
	Header [2]byte `struct:"[2]uint8"`
	Series uint8

	Trailer uint8 `struct:"uint8:3"`
	Channel uint8 `struct:"uint8:2"`
	Command uint8 `struct:"uint8:3"`
}

var series39Lookup = []struct {
	Channel
	Action
	b1 byte
	b2 byte
}{
	{Channel1, ActionWhite, Series39, 0x90},
	{Channel1, ActionWarm, Series39, 0x91},
	{Channel1, ActionSleep, Series39, 0xa1},

	{Channel2, ActionWhite, Series39, 0x94},
	{Channel2, ActionWarm, Series39, 0x95},
	{Channel2, ActionSleep, Series39, 0xaa},

	{Channel3, ActionWhite, Series39, 0x98},
	{Channel3, ActionWarm, Series39, 0x99},
	{Channel3, ActionSleep, Series39, 0xb3},
}

func (p Command) Encode() ([]byte, error) {
	var encoded []byte
	var err error
	if series := p.Action.series(); series == Series09 {
		// This almost looks structured.

		packet := Packet09{
			Header:  [2]byte{0x2c, 0x52},
			Series:  byte(series),
			Command: byte(p.Action),
			Channel: byte(p.Channel),
			Trailer: 1,
		}
		encoded, err = restruct.Pack(binary.LittleEndian, &packet)
		if err != nil {
			return nil, err
		}
		encoded = append(encoded, encoded[2]^encoded[3])
	} else {
		// This doesn't have an obvious structure.

		found := false
		var b1, b2 byte
		for _, e := range series39Lookup {
			if e.Channel == p.Channel && e.Action == p.Action {
				found = true
				b1 = e.b1
				b2 = e.b2
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("couldn't find encode for command %v", p)
		}
		encoded = []byte{0x2c, 0x52, b1, b2, b1 ^ b2}
	}

	return encoded, nil
}
