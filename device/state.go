package device

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
)

type State interface {
	Bindings() []Binding
}

type state struct {
	StateFile string
	State     State
}

func (state *state) LoadState() error {
	contents, err := ioutil.ReadFile(state.StateFile)

	if err != nil {
		return err
	}

	statePtr := reflect.New(reflect.TypeOf(state.State))
	statePtr.Elem().Set(reflect.ValueOf(state.State))

	if err := json.Unmarshal(contents, statePtr.Interface()); err != nil {
		return err
	}

	state.State = statePtr.Elem().Interface().(State)

	return nil
}

func (state *state) WriteState() error {
	contents, err := json.Marshal(state.State)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(state.StateFile, contents, 0640)
}

func GenerateStateConfig(name string, prefix string, state State) map[string]string {
	result := map[string]string{}

	result["name"] = name
	result["~"] = prefix

	for _, binding := range state.Bindings() {
		for k, v := range binding.Config() {
			result[k] = v
		}
	}

	return result
}
