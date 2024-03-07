{ pkgs, ... }: 
pkgs.mkShell {
   packages = with pkgs; [ go_1_22 ];
   shellHook = ''
    export PATH=$PWD/scripts:$PATH
   '';
}
