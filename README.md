# RAG Pipeline

## Description

Simplified RAG Pipeline in Go and Python.

This project assumes you have [`Nix`](https://nix.dev/install-nix) installed and relies on `shell.nix` for dev environment. The environment can be entered with the `nix-shell` command. Python dependencies are managed directly through Nix.

For first run, `make all` will build and run the project. `make build` only needs to be ran if changes are made to the Go components. Subsequent execution simply requires `make run`.

See the rest of [`Makefile`](./Makefile) for other tasks, such as `docker-up` and `clean`.

### Components

- Solana Data Extractor (Datasource & Processing)
- RAG Pipeline Orchestrator
- Data Cleaning & Enrichment
- Neo4j Integration (and React Dashboard)

#### `.env` File

There is an `.env.sample` file you can copy to `.env`. These environment variables are required for proper execution of this program and docker-compose configuration.