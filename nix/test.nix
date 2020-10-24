import (<nixpkgs> + "/nixos/tests/make-test-python.nix") ({ pkgs, lib, ... }:
{
  machine = { options, pkgs, ... }: {
    imports = [ ./module.nix ];

    services.hass_ir_adapter = {
      enable = true;
      config = builtins.readFile ../example.yaml;
    };

    environment.systemPackages = with pkgs; [ mosquitto jq ];
    services.mosquitto.enable = true;
    services.mosquitto.allowAnonymous = true;
    services.mosquitto.users.hass_ir_adapter = {
      acl = ["topic readwrite #"];
    };
    services.mosquitto.aclExtraConf = ''
      topic readwrite #
    '';
  };

  testScript = ''
    start_all()

    machine.wait_for_unit("mosquitto.service")
    machine.wait_for_unit("hass_ir_adapter.service")

    # Loggers to enhance test output
    machine.execute(
        'mosquitto_sub -t "homeassistant/#" -F "%t -> %p" | logger -t mqtt-homeassistant & disown'
    )
    machine.execute('mosquitto_sub -t "ir/#" -F "%t -> %x" | logger -t mqtt-ir & disown')

    # Command captures
    machine.execute("mosquitto_sub -t ir/ESP_1/send -F %x | tee /tmp/esp1.log &")
    machine.execute(
        "mosquitto_sub -t ir/tasmota/cmnd/IRhvac -F %p | tee /tmp/tasmota.log &"
    )
    machine.execute(
        "mosquitto_sub -t homeassistant/climate/living_room_ac/state -F %p | tee /tmp/climate.log &"
    )

    # Send state changes
    machine.execute(
        "mosquitto_pub -t homeassistant/climate/living_room_ac/set_mode -m cool"
    )
    machine.execute("mosquitto_pub -t homeassistant/climate/tasmota_ac/set_mode -m cool")

    machine.sleep(5)

    # Test state changes resulted in expected commands
    machine.succeed(
        'set -x && test "$(cat /tmp/esp1.log)" == 010ff123cb26010024030d000000000049'
    )
    machine.succeed(
        'set -x && test "$(jq -r --slurp "last | .Vendor" < /tmp/tasmota.log)" == HITACHI_AC424'
    )
    machine.succeed(
        'set -x && test "$(jq -r --slurp "last | .mode" < /tmp/climate.log)" == cool'
    )
  '';
})
