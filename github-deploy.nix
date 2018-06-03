{ lib, buildGoPackage, source ? lib.cleanSource ./. }:
let
  versionFile = builtins.readFile "${source}/version.go";
  versionMatch = builtins.match ".*\"([0-9]+\\.[0-9]+\\.[0-9]+)\".*" versionFile;
  version = builtins.head versionMatch;
in
buildGoPackage rec {
  name = "github-deploy-${version}";
  goPackagePath = "github.com/zimbatm/github-deploy";
  src = source;

  # FIXME: for some reason the go reference rewrites are failing
  allowGoReference = true;
}
