{ pkgs ? import <nixpkgs> {} }:

# Python 3.12
let
  python = pkgs.python312;
  pythonPackages = python.pkgs;
in

pkgs.mkShell {
  buildInputs = with pkgs; [
    # Go
    go
    gopls
    golangci-lint
    gotools
    go-tools
    golines

    # Python
    python
    pythonPackages.pip
    pythonPackages.numpy
    pythonPackages.pandas
    pythonPackages.sqlalchemy
    pythonPackages.psycopg2
    pythonPackages.neo4j
    pythonPackages.black  # Python formatter

    # Tools
    gnumake
    protobuf
    grpcurl

  ];

  shellHook = ''
    echo "RAG Pipeline development environment"
    echo "Go version: $(go version)"
    echo "Python version: $(python --version)"
    
    echo "Environment is ready!"
  '';
}