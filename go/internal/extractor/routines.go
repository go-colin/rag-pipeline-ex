package extractor

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// Gets the largest holders (token accounts) of the token then finds their Owners (wallets)
// Returns new WalletData
func (e *SolanaExtractor) InitWalletDataFromToken(
	tokenPubKey solana.PublicKey,
) ([]*WalletData, error) {
	// get top 20 holders for the token
	tokenAccounts, err := e.client.GetTokenLargestAccounts(
		context.TODO(),
		tokenPubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return nil, fmt.Errorf("getTokenLargestAccounts: %w", err)
	}

	// map holders from tokenAccounts
	extractHolders := func(accounts []*rpc.TokenLargestAccountsResult) []solana.PublicKey {
		holders := make([]solana.PublicKey, len(tokenAccounts.Value), len(tokenAccounts.Value))
		for i, acc := range accounts {
			holders[i] = acc.Address
		}
		return holders
	}
	holders := extractHolders(tokenAccounts.Value)

	// find owner wallets for list of holders
	// api limit has a range of 5
	allOwnerAccounts := make([]*rpc.Account, 0, len(holders))
	batchSize := 5
	for i := 0; i < len(holders); i += batchSize {
		// chunk holders into 5 at a time
		end := i + batchSize
		if end > len(holders) {
			end = len(holders)
		}

		ownerAccounts, err := e.client.GetMultipleAccounts(
			context.TODO(),
			holders[i:end]...,
		)
		if err != nil {
			return nil, fmt.Errorf("getMultipleAccounts: %w", err)
		}

		allOwnerAccounts = append(allOwnerAccounts, ownerAccounts.Value...)
	}

	// map holders from tokenAccounts
	extractOwners := func(accounts []*rpc.Account) ([]solana.PublicKey, error) {
		owners := make([]solana.PublicKey, len(tokenAccounts.Value), len(tokenAccounts.Value))
		for i, acc := range accounts {
			marshal, err := acc.Data.MarshalJSON()
			if err != nil {
				return nil, err
			}
			// owner is in value.data.parsed.owner
			fmt.Printf("%s", marshal)
			owners[i] = acc.Owner
		}
		return owners, nil
	}
	owners, err := extractOwners(allOwnerAccounts)
	if err != nil {
		return nil, fmt.Errorf("extractOwners: %w", err)
	}

	if e.litter != nil {
		e.litter.Dump(holders, owners)
	}

	return nil, nil
}

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

func (e *SolanaExtractor) GetProgramAccounts(walletData []*WalletData) error {

	for _, wallet := range walletData {
		fmt.Println(wallet)

		programAccounts, err := e.client.GetTokenAccountsByOwner(
			context.TODO(),
			wallet.Address,
			&rpc.GetTokenAccountsConfig{
				ProgramId: solana.TokenProgramID.ToPointer(),
			},
			nil,
		)
		if err != nil {
			return err
		}
		if e.litter != nil {
			fmt.Println("GetProgramAccounts Dump:")
			e.litter.Dump(programAccounts)
		}
	}

	return nil
}
