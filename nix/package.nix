{ lib, buildGo111Module ? null, buildGo112Module ? null }:

let
  goVersion =
    if ! isNull buildGo112Module then "1.12"
    else if ! isNull buildGo111Module then "1.11"
    else throw "Unknown version of Go";

  modSha256 = {
    "1.11" = "1zj90zvbn258sykl1kdh57y1y34vis0ygjjaicvs60m3v24ig2wf";
    "1.12" = "02ppjvaz5lkz0afwbcrf8xk986bs88dqzm0aphbkq0x9jdv7dhxr";
  }.${goVersion} or (throw "Missing modSha256 for go version ${goVersion}");

  buildGoModule = {
    "1.11" = buildGo111Module;
    "1.12" = buildGo112Module;
  }.${goVersion} or (throw "Missing buildGoModule for go version ${goVersion}");
in

buildGoModule {
  name = "hass_ir_adapter";
  src = lib.sourceFilesBySuffices ./.. [ ".go" ".mod" ".sum" ];
  inherit modSha256;
  subPackages = [ "." ];
  doCheck = true;
}
