package config

import (
	"fmt"
	"os"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
)

type Config struct {
	// General config
	DoLitter bool // if we should dump litter for debug

	// Solana config
	SolanaRPCURL string

	// RabbitMQ config
	RabbitMQURL string

	// Db config
	PostgresURL string
}

// Load config
func Load() (*Config, error) {
	// load .env
	godotenv.Load()

	cfg := &Config{
		DoLitter:     getEnvAsBool("APP_DEBUG", false),
		SolanaRPCURL: getEnv("SOLANA_RPC_URL", rpc.MainNetBeta_RPC),
		RabbitMQURL:  getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		PostgresURL:  getEnv("POSTGRES_URL", "postgres://user:password@localhost:5432/ragpipe"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil

}

// Validate config
func (c *Config) Validate() error {
	if c.SolanaRPCURL == "" {
		return fmt.Errorf("SOLANA_RPC_URL is required")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		return value == "true" || value == "1"
	}
	return fallback
}
