package types

type Entity struct {
	Id   string
	Type string
}

type Device struct {
	Entity  `yaml:",inline"`
	Emitter string
	Name    string
}

type Emitter struct {
	Entity `yaml:",inline"`
}

type Aircon struct {
	Device           `yaml:",inline"`
	TemperatureTopic string `yaml:"temperature_topic"`
}

type Light struct {
	Device `yaml:",inline"`
}
