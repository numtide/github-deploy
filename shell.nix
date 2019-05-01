let pkgs = import <nixpkgs> {}; in with pkgs;
mkShell {
  buildInputs = [
    go
  ];
  shellHook = ''
    export GO111MODULE=on
  '';
}
