{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/release-23.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachSystem [ "x86_64-linux" ] (system:
      let
        pkgs = import nixpkgs { inherit system; };
        github-deploy = import self {
          inherit system;
          inputs = null;
          nixpkgs = nixpkgs.legacyPackages.${system};
        };
        name = "github-deploy";
      in
      with pkgs;
      {
        devShell = github-deploy.devShell;
        packages.${name} = github-deploy.github-deploy;
        defaultPackage = github-deploy.defaultPackage;
      }
    );
}

