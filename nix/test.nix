import (<nixpkgs> + "/nixos/tests/make-test.nix") ({ pkgs, lib, ... }:
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
    startAll;

    $machine->waitForUnit('mosquitto.service');
    $machine->waitForUnit('hass_ir_adapter.service');

    $machine->execute('mosquitto_sub -t "homeassistant/#" -F "%t -> %p" | logger -t mqtt-homeassistant & disown');
    $machine->execute('mosquitto_sub -t "ir/#" -F "%t -> %x" | logger -t mqtt-ir & disown');

    $machine->execute('mosquitto_sub -t ir/ESP_1/send -F %x | tee /tmp/esp1.log &');
    $machine->execute('mosquitto_sub -t homeassistant/climate/living_room/state -F %p | tee /tmp/climate.log &');

    $machine->execute('mosquitto_pub -t homeassistant/climate/living_room/mode_command -m cool');

    $machine->sleep(5);

    $machine->succeed('set -x && test "$(cat /tmp/esp1.log)" == 010ff123cb26010024030d000000000049');
    $machine->succeed('set -x && test "$(jq -r --slurp "last | .mode" < /tmp/climate.log)" == cool');
  '';
})
