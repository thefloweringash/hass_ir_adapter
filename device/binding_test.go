package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testAutoState struct {
	TestInt      int     `json:"test_int_json" hass:"test_int_hass"`
	TestString   string  `json:"test_string_json" hass:"test_string_hass"`
	TestFloat    float32 `json:"test_float_json" hass:"test_float_hass"`
	NonHassField int
}

func (s testAutoState) Bindings() []Binding {
	return AutomaticBindings(s, "state")
}

func TestGenerateStateConfig(t *testing.T) {
	assert.Equal(t,
		map[string]string{
			"name":                            "my_name",
			"~":                               "my_prefix",
			"test_int_hass_command_topic":     "~/set_test_int_json",
			"test_int_hass_state_topic":       "~/state",
			"test_int_hass_state_template":    "{{ value_json.test_int_json }}",
			"test_string_hass_command_topic":  "~/set_test_string_json",
			"test_string_hass_state_topic":    "~/state",
			"test_string_hass_state_template": "{{ value_json.test_string_json }}",
			"test_float_hass_command_topic":   "~/set_test_float_json",
			"test_float_hass_state_topic":     "~/state",
			"test_float_hass_state_template":  "{{ value_json.test_float_json }}",
		},
		GenerateStateConfig("my_name", "my_prefix", testAutoState{}))
}

func TestBindings(t *testing.T) {
	bindings := AutomaticBindings(testAutoState{}, "state")

	tests := []struct {
		newState testAutoState
		binding  Binding
		value    string
	}{
		{newState: testAutoState{TestInt: 42},
			binding: bindings[0],
			value:   "42"},
		{newState: testAutoState{TestString: "forty-two"},
			binding: bindings[1],
			value:   "forty-two"},
		{newState: testAutoState{TestFloat: 42},
			binding: bindings[2],
			value:   "42.0"},
	}

	for _, test := range tests {
		newState, err := deriveState(testAutoState{}, test.binding.(InferredBinding), test.value)
		assert.Nil(t, err)
		assert.Equal(t, test.newState, newState)
	}
}
