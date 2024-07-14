package extractor

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type WalletData struct {
	Address           solana.PublicKey
	Transactions      []*rpc.TransactionSignature
	TokenBalances     map[solana.PublicKey]float64
	TxFrequency       int
	ProvidedLiquidity bool
}

// Make new wallet data from initial token and largest accounts values
func NewWalletDataFromTokenValues(
	token solana.PublicKey,
	values []*rpc.TokenLargestAccountsResult,
) []*WalletData {

	var walletData []*WalletData
	for _, account := range values {
		// make new wallet data
		wallet := &WalletData{
			Address: account.Address,
			// load initial token amount
			TokenBalances: map[solana.PublicKey]float64{
				token: *account.UiAmount,
			},
		}
		walletData = append(walletData, wallet)
	}

	return walletData
}
