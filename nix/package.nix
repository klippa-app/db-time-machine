{ pkgs, ... }:
pkgs.buildGoModule {
   name = "db-time-machine";
   version = "1.0.0";
   src = ../.;
   vendorHash = null;
}
