{ lib, buildGo111Module }:

# Also builds with buildGo112Module, but with a different modSha256.
buildGo111Module {
  name = "hass_ir_adapter";
  src = lib.sourceFilesBySuffices ./.. [ ".go" ".mod" ".sum" ];
  modSha256 = "1zj90zvbn258sykl1kdh57y1y34vis0ygjjaicvs60m3v24ig2wf";
  subPackages = [ "." ];
  doCheck = true;
}
