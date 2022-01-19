{ system ? builtins.currentSystem
, inputs ? import ./flake.lock.nix { }
, nixpkgs ? import inputs.nixpkgs {
    inherit system;
    # Makes the config pure as well. See <nixpkgs>/top-level/impure.nix:
    config = { };
    overlays = [ ];
  }
, buildGoPackage ? nixpkgs.buildGoPackage
}:
let
  versionFile = builtins.readFile ./version.go;
  versionMatch = builtins.match ".*\"([0-9]+\\.[0-9]+\\.[0-9]+)\".*" versionFile;
  version = builtins.head versionMatch;
  github-deploy = buildGoPackage
    rec {
      name = "github-deploy-${version}";
      goPackagePath = "github.com/zimbatm/github-deploy";
      src = nixpkgs.lib.cleanSource ./.;
      # FIXME: for some reason the go reference rewrites are failing
      allowGoReference = true;
    };
in
{
  inherit github-deploy;
  defaultPackage = github-deploy;
  devShell = nixpkgs.mkShell {
    buildInputs = with nixpkgs; [
      go
      gofumpt
      golangci-lint
      gopls
      nixpkgs-fmt
      just
      treefmt
    ];
  };
}
