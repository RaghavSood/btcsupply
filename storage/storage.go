package storage

import "github.com/RaghavSood/btcsupply/types"

type Storage interface {
	GetRecentLosses(limit int) ([]types.Loss, error)
	GetBlockLosses(hash string) ([]types.Loss, error)
	GetBlock(hash string) (types.Block, error)
}
