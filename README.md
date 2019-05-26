# Home-Assistant Compatible Adapter for MQTT Based IR blasters

Connects to MQTT and presents a [climate.mqtt][] compatible
device. Climate requests are translated into IR signals and emitted
via MQTT (again).

Intended for use with [thefloweringash/irsender][], but other backends
can be added.

For configuration examples see example.yaml.

[climate.mqtt]: https://www.home-assistant.io/components/climate.mqtt/

[thefloweringash/irsender]: [https://github.com/thefloweringash/irsender]
