{ lib, buildGo111Module }:

# Also builds with buildGo112Module, but with a different modSha256.
buildGo111Module {
  name = "hass_ir_adapter";
  src = lib.sourceFilesBySuffices ./.. [ ".go" ".mod" ".sum" ];
  modSha256 = "01621g84a93cqsc5wb3cjw2vl1r3x5l7kj2iz7jfzknyv73jaihc";
  doCheck = true;
}
