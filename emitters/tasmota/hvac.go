package tasmota

import (
	"encoding/json"
)

type StringBool bool

func (s StringBool) MarshalJSON() ([]byte, error) {
	if s {
		return json.Marshal("On")
	} else {
		return json.Marshal("Off")
	}
}

type Mode string
type FanSpeed string
type SwingV string
type SwingH string

const (
	ModeOff  Mode = "Off"
	ModeAuto Mode = "Auto"
	ModeCool Mode = "Cool"
	ModeHeat Mode = "Heat"
	ModeDry  Mode = "Dry"
	ModeFan  Mode = "Fan"

	FanSpeedAuto    FanSpeed = "Auto"
	FanSpeedOff     FanSpeed = "Off"
	FanSpeedMin     FanSpeed = "Min"
	FanSpeedLow     FanSpeed = "Low"
	FanSpeedMid     FanSpeed = "Mid"
	FanSpeedHigh    FanSpeed = "High"
	FanSpeedHighest FanSpeed = "Highest"

	SwingVAuto    SwingV = "Auto"
	SwingVOff     SwingV = "Off"
	SwingVMin     SwingV = "Min"
	SwingVLow     SwingV = "Low"
	SwingVMid     SwingV = "Mid"
	SwingVHigh    SwingV = "High"
	SwingVHighest SwingV = "Highest"

	SwingHAuto     SwingH = "Auto"
	SwingHOff      SwingH = "Off"
	SwingHLeftMax  SwingH = "LeftMax"
	SwingHLeft     SwingH = "Left"
	SwingHMid      SwingH = "Mid"
	SwingHRight    SwingH = "Right"
	SwingHRightMax SwingH = "RightMax"
	SwingHWide     SwingH = "Wide"
)

type TasmotaHvac struct {
	Vendor   string
	Power    StringBool
	Mode     Mode
	FanSpeed FanSpeed
	SwingV   SwingV
	SwingH   SwingH
	Celsius  StringBool
	Temp     float32
	Quiet    StringBool
	Turbo    StringBool
	Econo    StringBool
	Light    StringBool
	Filter   StringBool
	Clean    StringBool
	Beep     StringBool
	Sleep    int
}

func (self *TasmotaEmitter) EmitHVAC(command TasmotaHvac) error {
	// fill defaults
	if command.Mode == "" {
		command.Mode = ModeOff
	}

	if command.FanSpeed == "" {
		command.FanSpeed = FanSpeedAuto
	}

	if command.SwingV == "" {
		command.SwingV = SwingVOff
	}

	if command.SwingH == "" {
		command.SwingH = SwingHOff
	}

	jsonString, err := json.Marshal(&command)
	if err != nil {
		return err
	}

	self.logger.Printf("emitting with payload %v\n", string(jsonString))

	token := self.client.Publish(self.topic+"/IRhvac", 0, false, jsonString)
	token.Wait()
	return token.Error()
}
