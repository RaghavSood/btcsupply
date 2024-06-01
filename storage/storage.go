package storage

type Storage interface {
	GetRecentLosses(limit int) ([]types.Loss, error)
}
