package types

import "time"

type Loss struct {
	ID          int
	TxID        string
	BlockHash   string
	BlockHeight int64
	Vout        int
	Amount      *BigInt
	CreatedAt   time.Time
	BurnScript  string
}
