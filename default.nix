{ pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
, buildGoApplication ? pkgs.buildGoApplication
}:

buildGoApplication {
  nativeBuildInputs = [ pkgs.pkg-config ];
  buildInputs = [ pkgs.nixVersions.nix_2_23 ];
  pname = "nixdockertag";
  version = "1.0";
  src = ./.;
  pwd = ./.;
  go = pkgs.go_1_21;
  modules = ./gomod2nix.toml;
}
