package device

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Binding interface {
	RelativeTopic() string
	Config() map[string]string
	Apply(state State, value string) (State, error)
}

type stateField struct {
	Name             string
	JsonName         string
	HassName         string
	ValueTemplateKey string
	reflection       reflect.StructField
}

func (f stateField) CommandRelativeTopic() string {
	return fmt.Sprintf("set_%s", f.JsonName)
}

func (f stateField) CommandTopic(prefix string) string {
	return fmt.Sprintf("%s/%s", prefix, f.CommandRelativeTopic())
}

type BindingHook = func(state State, update func(State) (State, error)) (State, error)

type FieldBinding struct {
	Field          stateField
	TemplateSuffix string
	Hook           BindingHook
}

func (b FieldBinding) String() string {
	return fmt.Sprintf("bound_field(%s)", b.Field.reflection.Name)
}

func (b FieldBinding) Apply(state State, value string) (State, error) {
	update := func(state State) (State, error) {
		return deriveState(state, b, value)
	}
	return b.Hook(state, update)
}

func (b FieldBinding) RelativeTopic() string {
	return b.Field.CommandRelativeTopic()
}

func (b FieldBinding) Config() map[string]string {
	template := b.Field.ValueTemplateKey
	if template == "" {
		template = fmt.Sprintf("%s_%s_template", b.Field.HassName, b.TemplateSuffix)
	}
	return map[string]string{
		b.Field.HassName + "_command_topic": b.Field.CommandTopic("~"),
		b.Field.HassName + "_state_topic":   "~/state",
		template:                            fmt.Sprintf("{{ value_json.%s }}", b.Field.JsonName),
	}
}

const (
	TemplateKey = "template_key:"
)

func reflectOnState(state interface{}) []stateField {
	reflection := reflect.TypeOf(state)

	fields := make([]stateField, 0, reflection.NumField())

	for i := 0; i < reflection.NumField(); i++ {
		field := reflection.Field(i)
		hassTag, ok := field.Tag.Lookup("hass")
		if !ok {
			continue
		}

		jsonTag, ok := field.Tag.Lookup("json")
		if !ok {
			continue
		}
		jsonOpts := strings.SplitN(jsonTag, ",", 2)
		jsonName := jsonOpts[0]

		var valueTemplateKey string

		hassOpts := strings.Split(hassTag, ",")
		hassName := hassOpts[0]
		for _, hassOpt := range hassOpts {
			if strings.HasPrefix(hassOpt, TemplateKey) {
				valueTemplateKey = hassOpt[len(TemplateKey):]
			}
		}

		fields = append(fields, stateField{
			Name:             field.Name,
			JsonName:         jsonName,
			HassName:         hassName,
			ValueTemplateKey: valueTemplateKey,
			reflection:       field,
		})
	}

	if len(fields) == 0 {
		panic(fmt.Sprintf("ReflectOnState: empty reflection: %v", state))
	}

	return fields
}

type AutomaticBindingOptions struct {
	TemplateSuffix string
	UpdateHooks    map[string]BindingHook
}

func defaultHook(state State, update func(State) (State, error)) (State, error) {
	return update(state)
}

func AutomaticBindings(
	state State,
	options AutomaticBindingOptions,
) []Binding {
	if options.TemplateSuffix == "" {
		options.TemplateSuffix = "state"
	}

	fields := reflectOnState(state)
	bindings := make([]Binding, 0, len(fields))
	for _, field := range fields {
		hook, hasHook := options.UpdateHooks[field.Name]
		if !hasHook {
			hook = defaultHook
		}
		bindings = append(bindings, FieldBinding{
			Field:          field,
			TemplateSuffix: options.TemplateSuffix,
			Hook:           hook,
		})
	}
	return bindings
}

func deriveState(state State, binding FieldBinding, value string) (State, error) {
	stateReflection := reflect.TypeOf(state)
	newStatePtr := reflect.New(stateReflection)
	newStatePtr.Elem().Set(reflect.ValueOf(state))

	target := newStatePtr.Elem().FieldByIndex(binding.Field.reflection.Index)

	switch kind := target.Kind(); kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16:
		val, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, err
		}
		if target.OverflowUint(val) {
			return nil, fmt.Errorf("uint value overflow %v", val)
		}
		target.SetUint(val)

	case reflect.Int, reflect.Int8:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		if target.OverflowInt(val) {
			return nil, fmt.Errorf("int value overflow %v", val)
		}
		target.SetInt(val)

	case reflect.Float32:
		val, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, err
		}
		if target.OverflowFloat(val) {
			return nil, fmt.Errorf("float value overflow: %v", val)
		}
		target.SetFloat(val)

	case reflect.String:
		target.SetString(value)

	default:
		panic(fmt.Errorf("don't know how to set kind %s", kind))
	}

	return newStatePtr.Elem().Interface().(State), nil
}

type CallbackBinding struct {
	Topic    string
	Conf     map[string]string
	Callback func(state State, value string) (State, error)
}

func (b CallbackBinding) String() string {
	return fmt.Sprintf("callback(%s)", b.Topic)
}

func (b CallbackBinding) RelativeTopic() string {
	return b.Topic
}

func (b CallbackBinding) Config() map[string]string {
	return b.Conf
}

func (b CallbackBinding) Apply(state State, value string) (State, error) {
	return b.Callback(state, value)
}
