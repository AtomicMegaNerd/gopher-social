{
  description = "A udemy course for backend development with Go";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        # See https://github.com/golang-migrate/migrate/issues/1279
        go-migrate-pkg = pkgs.go-migrate.overrideAttrs (oldAttrs: {
          tags = [ "postgres" ];
        });
      in
      {
        devShell = pkgs.mkShell {
          # The packages we need for this project
          buildInputs = [
            pkgs.go_1_25
            pkgs.go-tools
            pkgs.gopls
            pkgs.golangci-lint
            pkgs.go-task
            go-migrate-pkg
            pkgs.air
            pkgs.rainfrog
            pkgs.posting
            pkgs.bash-language-server
            pkgs.docker-language-server
            pkgs.yaml-language-server
          ];
        };
      }
    );
}
