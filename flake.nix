{
  description = "A udemy course for backend development with Go";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    {
      self,
      nixpkgs,
    }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-darwin"
      ];
      buildPkgsConf =
        system:
        import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };
    in
    {
      devShells = nixpkgs.lib.genAttrs systems (
        system:
        let
          pkgs = buildPkgsConf system;
        in
        {
          default = pkgs.mkShell {
            # The packages we need for this project
            buildInputs = [
              pkgs.go_1_26
              pkgs.go-tools
              pkgs.gopls
              pkgs.golangci-lint
              pkgs.go-task
              pkgs.go-migrate
              pkgs.air
              pkgs.rainfrog
              pkgs.posting
              pkgs.bash-language-server
              pkgs.docker-language-server
              pkgs.yaml-language-server
              pkgs.efm-langserver
            ];
          };
        }
      );
    };
}
