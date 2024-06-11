package types

type TransactionSummary struct {
	Txid        string  `json:"tx_id"`
	Coinbase    bool    `json:"coinbase"`
	TotalLoss   *BigInt `json:"total_loss"`
	BlockHeight int64   `json:"block_height"`
	BlockHash   string  `json:"block_hash"`
}
