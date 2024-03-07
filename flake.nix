{
  description = "nix is love, nix is life";
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  inputs.flake-utils = {
    url = "github:numtide/flake-utils";
    inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, flake-utils, nixpkgs }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; };
      in rec {
        devShells.default = pkgs.mkShell {
           packages = with pkgs; [ go_1_22 ];
           shellHook = ''
            export PATH=$PWD/scripts:$PATH
           '';
        };
        packages.default = pkgs.callPackage ./nix/package.nix {};
        apps.default = flake-utils.lib.mkApp {
            drv = packages.deafault;
            exePath = /bin/dbtm;
        };
      }) // {
        overlays.default = (final: prev: rec {
          db-time-machine = self.packages."${final.system}".default;
        });
      };
}
