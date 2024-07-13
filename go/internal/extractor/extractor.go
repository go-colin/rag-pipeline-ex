package extractor

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sanity-io/litter"

	"github.com/go-colin/rag-pipeline-ex/internal/config"
)

type SolanaExtractor struct {
	client       *rpc.Client
	rabbitmqConn *amqp.Connection
	rabbitmqChan *amqp.Channel
	litter       bool // if true, dump litter for debug
}

func NewSolanaExtractor(cfg *config.Config) (*SolanaExtractor, error) {
	// Init client
	client := rpc.New(cfg.SolanaRPCURL)

	// Init RMQ
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &SolanaExtractor{
		client:       client,
		rabbitmqConn: conn,
		rabbitmqChan: ch,
		litter:       cfg.DoLitter,
	}, nil
}

func (e *SolanaExtractor) Run(tokenAddress string) error {
	// validate tokenAddress
	tokenPubKey, err := solana.PublicKeyFromBase58(tokenAddress)
	defer e.rabbitmqConn.Close()
	defer e.rabbitmqChan.Close()

	if err != nil {
		return err
	}
	walletData, err := e.ExtractData(tokenPubKey)
	if err != nil {
		return err
	}
	if err = e.PublishToRabbitMQ(walletData); err != nil {
		return err
	}

	return nil
}

// Extract data from Solana
func (e *SolanaExtractor) ExtractData(tokenPubKey solana.PublicKey) ([]WalletData, error) {

	out, err := e.client.GetTokenLargestAccounts(
		context.TODO(),
		tokenPubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return nil, err
	}

	if e.litter {
		litter.Dump(out)
	}
	return nil, nil
}

func (e *SolanaExtractor) PublishToRabbitMQ(data []WalletData) error {

	return nil
}
