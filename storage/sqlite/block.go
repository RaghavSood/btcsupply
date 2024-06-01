package sqlite

import (
	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetBlock(hash string) (types.Block, error) {
	var block types.Block
	err := d.db.QueryRow("SELECT * FROM blocks WHERE block_hash = ?", hash).Scan(&block.ID, &block.BlockHeight, &block.BlockHash, &block.BlockTimestamp, &block.ParentBlockHash, &block.NumTransactions, &block.BlockReward, &block.FeesReceived, &block.CreatedAt)
	if err != nil {
		return types.Block{}, err
	}

	return block, nil
}
