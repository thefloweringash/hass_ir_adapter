{ lib, buildGo111Module }:

# Also builds with buildGo112Module, but with a different modSha256.
buildGo111Module {
  name = "hass_ir_adapter";
  src = lib.sourceFilesBySuffices ./.. [ ".go" ".mod" ".sum" ];
  modSha256 = "1lvkl2ay0i21vwywnai8hl3nm1pb3hibpg2r4bq8608pg5km4kxa";
}
