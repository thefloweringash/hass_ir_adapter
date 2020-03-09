{ lib, buildGoModule }:

let
  goVersion = lib.versions.majorMinor
    (buildGoModule {src = ""; modSha256 = ""; pname = ""; version="";}).go.version;

  modSha256 = {
    "1.11" = "1zj90zvbn258sykl1kdh57y1y34vis0ygjjaicvs60m3v24ig2wf";
    "1.12" = "02ppjvaz5lkz0afwbcrf8xk986bs88dqzm0aphbkq0x9jdv7dhxr";
    "1.13" = "02ppjvaz5lkz0afwbcrf8xk986bs88dqzm0aphbkq0x9jdv7dhxr";
    "1.14" = "02ppjvaz5lkz0afwbcrf8xk986bs88dqzm0aphbkq0x9jdv7dhxr";
  }.${goVersion} or (throw "Missing modSha256 for go version ${goVersion}");
in

buildGoModule {
  name = "hass_ir_adapter";
  src = lib.sourceFilesBySuffices ./.. [ ".go" ".mod" ".sum" ];
  inherit modSha256;
  subPackages = [ "." ];
  doCheck = true;
}
