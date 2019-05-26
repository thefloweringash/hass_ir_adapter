{ config, lib, pkgs, ... }:

let
  cfg = config.services.hass_ir_adapter;

  configFile = pkgs.writeTextFile {
    name = "hass_ir_adapter_config.yaml";
    text = cfg.config;
  };
in

with lib;

{
  options = with types; {
    services.hass_ir_adapter = {
      enable = mkEnableOption "Home-Assistant IR Adapter";

      config = mkOption {
        type = str;
        description = "Literal config as YAML";
      };

      package = mkOption {
        type = package;
        default = pkgs.hass_ir_adapter;
        defaultText = "pkgs.hass_ir_adapter";
      };
    };
  };

  config = mkMerge [
    {
      nixpkgs.overlays = [ (import ./overlay.nix) ];
    }

    (mkIf cfg.enable {
      systemd.services.hass_ir_adapter = {
        wantedBy = [ "multi-user.target" ];
        after = [ "network.target" ];
        serviceConfig = {
          DynamicUser = true;
          StateDirectory = "hass_ir_adapter";
          Restart = "always";
          ExecStart = "${cfg.package}/bin/hass_ir_adapter " +
            "-config-file=${configFile} -state-dir=/var/lib/hass_ir_adapter";
        };
      };
    })
  ];
}
