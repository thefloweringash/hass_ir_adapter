{ lib, buildGoModule }:

buildGoModule {
  name = "hass_ir_adapter";
  src = lib.sourceFilesBySuffices ./.. [ ".go" ".mod" ".sum" ];
  vendorSha256 = "1v9cdj4py769dn7q7pmp9kfs6h08kab8sc2pfkrsc9ll2dnngxdj";
  subPackages = [ "." ];
  doCheck = true;
}
