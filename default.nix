{ pkgs ? import (import ./nix/sources.nix).nixpkgs { } }:
pkgs.callPackage ./github-deploy.nix { }
