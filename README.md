# Home-Assistant Compatible Adapter for MQTT Based IR blasters

Connects to MQTT and presents a Home Assistant [climate.mqtt][] compatible
device. Climate requests are translated into IR signals and sent via a hardware
emitter.

Supported devices:
 - Panasonic Lights
 - Daiko Lights
 - Mitsubishi AC with remote GP-82
 - Mitsubishi AC with remote RH-101
 - Generic [Tasmota IRhvac][]

Supported emitters:
 - [thefloweringash/irsender][] over mqtt
 - [Tasmota][] over mqtt
 - [IRKit][] over http

For configuration examples see example.yaml.

[climate.mqtt]: https://www.home-assistant.io/components/climate.mqtt/

[thefloweringash/irsender]: https://github.com/thefloweringash/irsender
[Tasmota]: https://tasmota.github.io/docs/
[Tasmota IRhvac]: https://tasmota.github.io/docs/Tasmota-IR/#sending-irhvac-commands
[IRKit]: https://getirkit.com/en/
