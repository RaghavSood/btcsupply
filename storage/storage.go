package storage

import "github.com/RaghavSood/btcsupply/types"

type Storage interface {
	GetRecentLosses(limit int) ([]types.Loss, error)
	GetBlockLosses(hash string) ([]types.Loss, error)
	GetTransactionLosses(hash string) ([]types.Loss, error)

	GetLossyBlocks(limit int) ([]types.Block, error)
	GetBlock(hash string) (types.Block, error)

	GetLossNote(noteID string) (types.LossNote, error)
	GetLossNotes(noteIDs []string) ([]types.LossNote, error)

	GetTransactionDetail(hash string) (types.TransactionDetail, error)
}
