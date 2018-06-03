let pkgs = import <nixpkgs> {}; in with pkgs;
mkShell {
  buildInputs = [
    go
    dep
  ];
  shellHook = ''
    export GOPATH=$HOME/go
  '';
}
