package types

import "time"

type TransactionQueue struct {
	ID          int64     `json:"id"`
	Txid        string    `json:"txid"`
	BlockHeight int64     `json:"block_height"`
	CreatedAt   time.Time `json:"created_at"`
}
