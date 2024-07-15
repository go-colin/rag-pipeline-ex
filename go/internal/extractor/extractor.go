package extractor

import (
	"fmt"
	"io"
	"reflect"

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
	litter       *litter.Options
}

func NewSolanaExtractor(cfg *config.Config) (*SolanaExtractor, error) {
	// Init client
	client := rpc.New(cfg.SolanaRPCURL)

	// Init RMQ
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}

	// litter config:
	var litterOptions *litter.Options
	if cfg.DoLitter {
		litterOptions = &litter.Options{
			//Compact: true,
			// Custom Dump Function
			DumpFunc: func(v reflect.Value, w io.Writer) bool {
				// call String() on solana.PublicKey types.
				if v.Type() == reflect.TypeOf(solana.PublicKey{}) {
					w.Write(
						[]byte(
							fmt.Sprintf(" '%s'", v.MethodByName("String").Call(nil)[0].String()),
						),
					)
					return true
				}
				return false
			},
		}
	}

	return &SolanaExtractor{
		client:       client,
		rabbitmqConn: conn,
		rabbitmqChan: ch,
		litter:       litterOptions,
	}, nil
}

func (e *SolanaExtractor) Run(tokenAddress string) error {
	// validate tokenAddress
	tokenPubKey, err := solana.PublicKeyFromBase58(tokenAddress)
	defer e.rabbitmqConn.Close()
	defer e.rabbitmqChan.Close()

	if err != nil {
		return fmt.Errorf("tokenAddress: %w", err)
	}

	// Extract Data from token address
	walletData, err := e.ExtractData(tokenPubKey)
	if err != nil {
		return fmt.Errorf("extractData: %w", err)
	}

	// Publish wallet data to Queue
	if err = e.PublishToRabbitMQ(walletData); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

// Extract data from Solana
func (e *SolanaExtractor) ExtractData(tokenPubKey solana.PublicKey) ([]*WalletData, error) {

	// make initial walletdata from fetched wallets for token
	walletData, err := e.InitWalletDataFromToken(tokenPubKey)
	if err != nil {
		return nil, fmt.Errorf("initWalletData: %w", err)
	}

	// Get recent transactions for wallets
	/* if err = e.GetTransactions(walletData); err != nil {
		return nil, err
	} */

	// Get token accounts for wallets
	// if err = e.GetProgramAccounts(walletData[0:2]); err != nil {
	// 	return nil, fmt.Errorf("getProgramAccounts: %w", err)
	// }

	return walletData, nil
}

func (e *SolanaExtractor) PublishToRabbitMQ(data []*WalletData) error {

	return nil
}
