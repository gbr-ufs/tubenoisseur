# SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
#
# SPDX-License-Identifier: GPL-3.0-or-later

{
  description = "YouTuber Reference Scraper.";
  inputs.nixpkgs.url = "github:nixos/nixpkgs/master";
  outputs =
    { nixpkgs, ... }:
    let
      inherit (nixpkgs) lib;
      forAllSystems = lib.genAttrs lib.systems.flakeExposed;
    in
      {
        devShells = forAllSystems (
          system:
          let
            pkgs = nixpkgs.legacyPackages.${system};
            emacsSettings = pkgs.writeText "dir-locals.el" ''
              ((nil . ((eval . (with-eval-after-load 'apheleia
                         (add-to-list 'apheleia-formatters
                           '(golines . ("golines" "--base-formatter=goimports" "-w" filepath)))))))
               (go-mode . ((apheleia-formatter . (golines))
                           (eglot-server-programs . ((go-mode . ("rass" "--" "gopls" "--" "golangci-lint-langserver" :initializationOptions (:command ["golangci-lint" "run" "--output.json.path" "stdout" "--show-stats=false" "--issues-exit-code=1"])))))))
               (go-ts-mode . ((apheleia-formatter . (golines))
                              (eglot-server-programs . ((go-ts-mode . ("rass" "--" "gopls" "--" "golangci-lint-langserver" :initializationOptions (:command ["golangci-lint" "run" "--output.json.path" "stdout" "--show-stats=false" "--issues-exit-code=1"]))))))))
            '';
            neovimSettings = pkgs.writeText "nvim.lua" ''
              vim.lsp.config("gopls", {
                cmd = { "gopls" }
              })
              vim.lsp.enable("gopls")
              vim.lsp.config("golangci-lint-langserver", {
                cmd = { "golangci-lint-langserver" },
                init_options = {
                  command = { "golangci-lint", "run", "--output.json.path", "stdout", "--show-stats=false", "--issues-exit-code=1" }
                }
              })
              vim.lsp.enable("golangci-lint-langserver")
            '';
            VSCodeSettings = pkgs.writeText "vscode-settings.json" ''
              {
                "editor.formatOnSave": true,
                "go.formatTool": "goimports",
                "go.alternateTools": {
                  "goimports": "golines"
                },
                "go.formatFlags": [
                  "--base-formatter=goimports",
                  "-w"
                ],
                "go.lintTool": "golangci-lint",
                  "go.lintOnSave": "package",
                    "[go]": {
                      "editor.defaultFormatter": "golang.go"
                    }
                }
            '';
            zedSettings = pkgs.writeText "zed-settings.json" ''
              {
                "format_on_save": "on",
                "languages": {
                  "Go": {
                    "formatter": {
                      "external": {
                        "command": "golines",
                        "arguments": [
                          "--base-formatter=goimports",
                          "-w"
                        ]
                      }
                    },
                    "language_servers": [ "gopls", "golangci-lint-langserver" ]
                  }
                },
                "lsp": {
                  "gopls": {
                    "binary": {
                      "path_lookup": true,
                      "path": "gopls"
                    }
                  },
                  "golangci-lint-langserver": {
                    "binary": {
                      "path_lookup": true,
                      "path": "golangci-lint-langserver"
                    },
                    "initialization_options": {
                      "command": [
                        "golangci-lint",
                        "run",
                        "--output.json.path",
                        "stdout",
                        "--show-stats=false",
                        "--issues-exit-code=1"
                      ]
                    }
                  }
                }
              }
            '';
          in
            {
              default = pkgs.mkShell {
                packages = with pkgs; [
                  cocogitto
                  delve
                  go
                  golangci-lint
                  golangci-lint-langserver
                  golines
                  gopls
                  goreleaser
                  gotools
                  govulncheck
                  hugo
                  nixd
                  rassumfrassum
                  reuse
                  vhs
                ];
                shellHook = ''
                  mkdir -p .vscode
                  mkdir -p .zed
                  ln -sf ${emacsSettings} .dir-locals.el
                  ln -sf ${neovimSettings} .nvim.lua
                  ln -sf ${VSCodeSettings} .vscode/settings.json
                  ln -sf ${zedSettings} .zed/settings.json
                '';
              };
            }
        );
        formatter = forAllSystems (system: nixpkgs.legacyPackages.${system}.nixfmt);
      };
}
