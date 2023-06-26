{ pkgs ? import <nixpkgs> {} }:

let
    oldPkgs = import (builtins.fetchTarball {
        url = "https://github.com/NixOS/nixpkgs/archive/79b3d4bcae8c7007c9fd51c279a8a67acfa73a2a.tar.gz";
    }) {};
in

pkgs.mkShell {
  nativeBuildInputs = with pkgs.buildPackages; [
    libwebp
    libaom
    protobuf
    protoc-gen-go
    pkg-config
    oldPkgs.python37Full
    gcc
    libxcrypt
    libtensorflow
    python310Packages.tensorflow-bin
  ];
}
