# Variables

# Paths
GO_SRC := ./go
PYTHON_SRC := ./python

.PHONY: all build clean run test docker-up docker-down

all: build run

# Build Go components
build-go:
	cd $(GO_SRC)/cmd/solana_extractor && go build -o ./bin/solana_extractor
	cd $(GO_SRC)/cmd/rag_orchestrator && go build -o ./bin/rag_orchestrator

# Build
build: build-go

# Clean build artifacts
clean:
	rm -f $(GO_SRC)/cmd/solana_extractor/bin/solana_extractor
	rm -f $(GO_SRC)/cmd/rag_orchestrator/bin/rag_orchestrator
#	find . -type d -name __pycache__ -exec rm -r {} +

# Run the entire pipeline
run: docker-up
	$(GO_SRC)/cmd/solana_extractor/bin/solana_extractor
#	python $(PYTHON_SRC)/data_cleaning/cleaner.py
#	python $(PYTHON_SRC)/data_enrichment/enricher.py
#	python $(PYTHON_SRC)/neo4j_integration/neo4j_client.py
#	$(GO_SRC)/cmd/rag_orchestrator/bin/rag_orchestrator

# Run tests
test:
	cd $(GO_SRC) && go test ./...
	cd $(PYTHON_SRC) && python -m unittest discover

# Start Docker services
docker-up:
	docker-compose up -d

# Stop Docker services
docker-down:
	docker-compose down