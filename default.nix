{ pkgs ? import <nixpkgs> {} }:
pkgs.callPackage ./github-deploy.nix {}
