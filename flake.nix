{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/release-24.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachSystem [ "x86_64-linux" ] (system:
      let
        github-deploy = import self {
          inherit system;
          inputs = null;
          nixpkgs = nixpkgs.legacyPackages.${system};
        };
        name = "github-deploy";
      in
      {
        devShell = github-deploy.devShell;
        packages.${name} = github-deploy.github-deploy;
        defaultPackage = github-deploy.defaultPackage;

        checks.${name} = github-deploy.github-deploy;
      }
    );
}

