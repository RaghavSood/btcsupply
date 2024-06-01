package sqlite

import (
	"database/sql"

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

func (d *SqliteBackend) GetLossyBlocks(limit int) ([]types.Block, error) {
	rows, err := d.db.Query("SELECT * FROM blocks ORDER BY block_height DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}

	blocks, err := scanBlocks(rows)
	if err != nil {
		return nil, err
	}

	return blocks, nil
}

func scanBlocks(rows *sql.Rows) ([]types.Block, error) {
	var blocks []types.Block
	for rows.Next() {
		var block types.Block
		err := rows.Scan(&block.ID, &block.BlockHeight, &block.BlockHash, &block.BlockTimestamp, &block.ParentBlockHash, &block.NumTransactions, &block.BlockReward, &block.FeesReceived, &block.CreatedAt)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}
