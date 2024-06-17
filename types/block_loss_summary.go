package types

import "time"

type BlockLossSummary struct {
	BlockHeight    int64
	BlockHash      string
	LossOutputs    int
	TotalLost      *BigInt
	BlockTimestamp time.Time
}
