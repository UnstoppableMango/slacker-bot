{
  description = "A Discord bot for the Slackers";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    systems.url = "github:nix-systems/default";
    flake-parts.url = "github:hercules-ci/flake-parts";

    apis = {
      url = "github:unmango/apis";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        systems.follows = "systems";
        flake-parts.follows = "flake-parts";
        treefmt-nix.follows = "treefmt-nix";
      };
    };

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

          version = "0.0.1";
          slacker-bot = buildGoApplication {
            pname = "slacker-bot";
            inherit version;

            go = pkgs.go_1_26;
            src = lib.cleanSource ./.;
            modules = ./gomod2nix.toml;
          };

          ctr = pkgs.dockerTools.streamLayeredImage {
            name = "slacker-bot";
            tag = version;

            contents = [
              pkgs.cacert
              (pkgs.buildEnv {
                name = "image-root";
                paths = [ slacker-bot ];
                pathsToLink = [ "/bin" ];
              })
            ];

            config = {
              Entrypoint = [ "/bin/slacker-bot" ];
            };
          };
        in
        {
          packages = {
            inherit slacker-bot ctr;
            default = slacker-bot;
          };

          apps = {
            slacker-bot = {
              type = "app";
              program = "${slacker-bot}/bin/slacker-bot";
            };
            gopls = {
              type = "app";
              program = "${pkgs.gopls}/bin/gopls";
            };
          };

          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              buf
              gnumake
              go_1_26
              gomod2nix
              gopls
              nixfmt
              podman
              podman-compose
              skopeo
              uutils-findutils
            ];

            BUF = "${pkgs.buf}/bin/buf";
            FIND = "${pkgs.uutils-findutils}/bin/find";
            GO = "${pkgs.go_1_26}/bin/go";
            GOMOD2NIX = "${gomod2nix}/bin/gomod2nix";
            PODMAN = "${pkgs.podman}/bin/podman";

            PODMAN_COMPOSE_WARNING_LOGS = "false";
          };

          treefmt.programs = {
            actionlint.enable = true;
            gofmt.enable = true;
            jsonfmt.enable = true;
            nixfmt.enable = true;
            yamllint.enable = true;
          };
        };
    };
}
