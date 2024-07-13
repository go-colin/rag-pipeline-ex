package extractor

type WalletData struct {
	Address           string
	Transactions      []string
	TokenBalances     map[string]uint64
	TxFrequency       int
	ProvidedLiquidity bool
}
