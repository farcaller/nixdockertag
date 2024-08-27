{
  description = "Nix-wrapped docker image tags";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
    let
      imageNames = builtins.attrNames (builtins.readDir ./images);
    in
    {
      images = builtins.listToAttrs (
        map
          (name: {
            name = builtins.head (builtins.split "\\." name);
            value = let imageInfo = import ./images/${name}; in {
              inherit imageInfo;
              refTag = "${imageInfo.image}:${imageInfo.followTag}";
              refHash = "${imageInfo.image}@sha256:${imageInfo.hash}";
            };
          })
          imageNames);
    } // flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.callPackage ./. {
          inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
        };
        devShells.default = pkgs.callPackage ./shell.nix {
          inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
        };
      });
}
