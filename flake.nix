{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    mk.url = "github:x0k/mk";
  };
  outputs =
    {
      self,
      nixpkgs,
      mk,
    }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.${system} = {
        default = pkgs.mkShell {
          buildInputs = [
            mk.packages.${system}.default
            pkgs.go_1_23
            pkgs.tinygo
            pkgs.air
            pkgs.go-migrate
            pkgs.golangci-lint
            pkgs.sqlc
          ];
        };
      };
    };
}
