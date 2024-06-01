package types

import "time"

type Block struct {
	ID              int
	BlockHeight     int
	BlockHash       string
	BlockTimestamp  time.Time
	ParentBlockHash string
	NumTransactions int
	BlockReward     BigInt
	FeesReceived    BigInt
	CreatedAt       time.Time
}
