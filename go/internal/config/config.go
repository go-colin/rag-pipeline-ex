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
		RabbitMQURL:  getEnv("RABBITMQ_URL", tryBuildRabbit()),
		PostgresURL:  getEnv("POSTGRES_URL", tryBuildPostgres()),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("cfg validate: %w", err)
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

// Attempts to build RabbitMQURL. Will return empty if User and Pass are not set in env.
func tryBuildRabbit() string {
	user := getEnv("RABBITMQ_DEFAULT_USER", "")
	if user == "" {
		return ""
	}
	pass := getEnv("RABBITMQ_DEFAULT_PASS", "")
	if pass == "" {
		return ""
	}
	// default port fallback
	port := getEnv("RABBITMQ_NODE_PORT", "5672")
	host := getEnv("RABBITMQ_NODE_IP_ADDRESS", "localhost")

	// build and return connection string
	return fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)
}

// Attempts to build PostgresURL. Will return empty if User, Pass and Db are not set in env.
func tryBuildPostgres() string {
	user := getEnv("POSTGRES_USER", "")
	if user == "" {
		return ""
	}
	pass := getEnv("POSTGRES_PASSWORD", "")
	if pass == "" {
		return ""
	}
	db := getEnv("POSTGRES_DB", "")
	if db == "" {
		return ""
	}
	// default port fallback
	port := getEnv("POSTGRES_PORT", "5432")
	host := getEnv("PGHOST", "localhost")

	// build and return connection string
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, db)
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
