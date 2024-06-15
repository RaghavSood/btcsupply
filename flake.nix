{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    systems.url = "github:nix-systems/default";
    devenv.url = "github:cachix/devenv/1e4701fb1f51f8e6fe3b0318fc2b80aed0761914";
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  outputs = { self, nixpkgs, devenv, systems, ... } @ inputs:
    let
      forEachSystem = nixpkgs.lib.genAttrs (import systems);
    in
    {
      devShells = forEachSystem
        (system:
          let
            pkgs = nixpkgs.legacyPackages.${system};
          in
          {
            default = devenv.lib.mkShell {
              inherit inputs pkgs;
              modules = [
                {
                  languages.go = {
                    enable = true;
                  };

                  packages = with pkgs; [
                    goose
                    flyctl
                    air
                    tailwindcss
                  ];

                  enterShell = ''
                    echo "btcsupply shell started!"
                  '';
                }
              ];
            };
          });
      packages = forEachSystem
        (system:
          let
            pkgs = nixpkgs.legacyPackages.${system};
          in
          {
            default = pkgs.buildGoModule {
              name = "btcsupply";

              src = ./.;
              vendorHash = "sha256-DRfI+5OUKwoIPYwPghE4lIclvXNAQr45d1FYFjicmp8=";

              buildInputs = with pkgs; [
                go
                tailwindcss
                git
              ];

              subPackages = [ "cmd/btcsupply" ];

              preBuild = ''
                substituteInPlace main.go --replace-fail tailwindcss ${pkgs.tailwindcss}/bin/tailwindcss
                go generate main.go
              '';

              doCheck = false;

            };
          });
    };
}
