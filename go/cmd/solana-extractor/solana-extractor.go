package main

import (
	"log"

	"github.com/go-colin/rag-pipeline-ex/internal/config"
	"github.com/go-colin/rag-pipeline-ex/internal/extractor"
)

func main() {

	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration %v", err)
	}

	// init extractor
	solEx, err := extractor.NewSolanaExtractor(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize extractor: %v", err)
	}

	// Run extractor
	if err = solEx.Run("3Rcc6tMyS7ZEa29dxV4g3J5StorS9J1dn98gd42pZTk1"); err != nil {
		log.Fatalf("Failed to run extractor: %v", err)
	}
}
