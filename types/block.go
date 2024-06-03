package types

import (
	"time"

	btypes "github.com/RaghavSood/btcsupply/bitcoinrpc/types"
)

type Block struct {
	ID              int
	BlockHeight     int64
	BlockHash       string
	BlockTimestamp  time.Time
	ParentBlockHash string
	NumTransactions int
	CreatedAt       time.Time
}

func FromRPCBlock(block btypes.Block) Block {
	return Block{
		BlockHeight:     block.Height,
		BlockHash:       block.Hash,
		BlockTimestamp:  time.Unix(int64(block.Time), 0),
		ParentBlockHash: block.Previousblockhash,
		NumTransactions: block.NTx,
	}
}
