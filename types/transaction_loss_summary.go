package types

type TransactionLossSummary struct {
	Txid        string  `json:"tx_id"`
	TotalLoss   *BigInt `json:"total_loss"`
	BlockHeight int64   `json:"block_height"`
	BlockHash   string  `json:"block_hash"`
}
