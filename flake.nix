{
  outputs = { self, nixpkgs, flake-utils }: flake-utils.lib.eachDefaultSystem (system: {
    devShells.default = with import nixpkgs { inherit system; }; mkShell {
      buildInputs = [ pkg-config libgit2_1_5 ];
    };
  });
}
