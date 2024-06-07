package types

type HeightLossSummary struct {
	BlockHeight int64   `json:"block_height"`
	TotalLoss   *BigInt `json:"total_loss"`
}
