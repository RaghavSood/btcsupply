package storage

import "github.com/RaghavSood/btcsupply/types"

type Storage interface {
	GetRecentLosses(limit int) ([]types.Loss, error)
}
