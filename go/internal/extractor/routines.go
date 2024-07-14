package extractor

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/sanity-io/litter"
)

// Get latest transaction signatures for provided wallets
func (e *SolanaExtractor) GetTransactions(walletData []*WalletData) error {

	limit := 20
	// iterate wallets and load transaction signatures
	for _, wallet := range walletData {
		transactions, err := e.client.GetSignaturesForAddressWithOpts(
			context.TODO(),
			wallet.Address,
			&rpc.GetSignaturesForAddressOpts{Limit: &limit},
		)
		if err != nil {
			return err
		}

		/* if e.litter {
			fmt.Println("GetTransactions Dump:")
			litter.Dump(transactions)
		} */

		wallet.Transactions = transactions
	}

	return nil
}

// Get token accounts for provided wallets
func (e *SolanaExtractor) GetTokenAccounts(walletData []*WalletData) error {

	//sliceLength := uint64(20)
	for _, wallet := range walletData {
		fmt.Println(wallet.Address)
		tokenAccounts, err := e.client.GetTokenAccountsByOwner(
			context.TODO(),
			wallet.Address,
			&rpc.GetTokenAccountsConfig{
				ProgramId: solana.TokenProgramID.ToPointer(),
			},
			&rpc.GetTokenAccountsOpts{
				Encoding: solana.EncodingBase64Zstd,
			},
		)
		if err != nil {
			return err
		}
		if e.litter {
			fmt.Println("GetTokenAccounts Dump:")
			litter.Dump(tokenAccounts)
		}
	}

	return nil
}
