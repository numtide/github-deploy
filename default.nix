{
  system ? builtins.currentSystem,
  inputs ? import ./flake.lock.nix { },
  nixpkgs ? import inputs.nixpkgs {
    inherit system;
    # Makes the config pure as well. See <nixpkgs>/top-level/impure.nix:
    config = { };
    overlays = [ ];
  },
  buildGoModule ? nixpkgs.buildGoModule,
}:
let
  inherit (nixpkgs) lib;
  versionFile = lib.readFile ./version.go;
  versionMatch = lib.match ".*\"([0-9]+\\.[0-9]+\\.[0-9]+)\".*" versionFile;
  version = lib.head versionMatch;
  github-deploy = buildGoModule {
    name = "github-deploy-${version}";
    src = nixpkgs.lib.cleanSource ./.;
    vendorHash = "sha256-vcm3UhqwiUCdope9TpfO/CQQsxM0eh8nhvOCv4bGCng=";
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
