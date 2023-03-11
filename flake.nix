{
  description = "Nix-wrapped docker image tags";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.poetry2nix = {
    url = "github:nix-community/poetry2nix";
    inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, poetry2nix, ... }:
    let
      imageNames = builtins.attrNames (builtins.readDir ./images);
    in
    {
      images = builtins.listToAttrs (
        map (name: {
          name = builtins.head (builtins.split "\\." name);
          value = import ./images/${name};
          }) imageNames);
    } // flake-utils.lib.eachDefaultSystem (system:
      let
        inherit (poetry2nix.legacyPackages.${system}) mkPoetryApplication;
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.updater = mkPoetryApplication {
          python = pkgs.python311;
          projectDir = ./.;
        };
      });
}
