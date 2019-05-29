{ nixpkgs ? <nixpkgs> }:

(import nixpkgs {}).callPackage ./nix/package.nix { }
