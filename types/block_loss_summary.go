package types

type BlockLossSummary struct {
	BlockHeight int64
	BlockHash   string
	LossOutputs int
	TotalLost   *BigInt
}
