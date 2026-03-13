{
  description = "A Discord bot for the Slackers";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    systems.url = "github:nix-systems/default";
    flake-parts.url = "github:hercules-ci/flake-parts";

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.inputs.systems.follows = "systems";
    };

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;
      imports = [ inputs.treefmt-nix.flakeModule ];

      perSystem =
        {
          inputs',
          pkgs,
          lib,
          ...
        }:
        let
          inherit (inputs'.gomod2nix.legacyPackages) gomod2nix buildGoApplication;

          slacker-bot = buildGoApplication {
            pname = "slacker-bot";
            version = "0.0.1";
            src = lib.cleanSource ./.;

            modules = ./gomod2nix.toml;
          };
        in
        {
          packages = {
            inherit slacker-bot;
            default = slacker-bot;
          };

          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              gnumake
              go
              gomod2nix
              nixfmt
              uutils-findutils
            ];

            FIND = "${pkgs.uutils-findutils}/bin/find";
            GO = "${pkgs.go}/bin/go";
            NIXFMT = "${pkgs.nixfmt}/bin/nixfmt";
          };

          treefmt = {
            programs.nixfmt.enable = true;
          };
        };
    };
}
