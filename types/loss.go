package types

import "time"

type Loss struct {
	ID        int
	TxID      string
	BlockHash string
	Vout      int
	Amount    *BigInt
	CreatedAt time.Time
}
