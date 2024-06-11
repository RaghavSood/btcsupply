package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetBlock(identifier string) (types.Block, error) {
	var block types.Block
	err := d.db.QueryRow("SELECT * FROM blocks WHERE block_hash = ? OR block_height = ?", identifier, identifier).Scan(&block.ID, &block.BlockHeight, &block.BlockHash, &block.BlockTimestamp, &block.ParentBlockHash, &block.NumTransactions, &block.CreatedAt)
	if err != nil {
		return types.Block{}, err
	}

	return block, nil
}

func (d *SqliteBackend) GetLossyBlocks(limit int) ([]types.BlockLossSummary, error) {
	query := `SELECT block_height, block_hash, COUNT(*), SUM(total_loss) FROM transaction_summary GROUP BY block_height ORDER BY block_height DESC LIMIT ?`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []types.BlockLossSummary
	for rows.Next() {
		var summary types.BlockLossSummary
		if err := rows.Scan(
			&summary.BlockHeight, &summary.BlockHash, &summary.LossOutputs, &summary.TotalLost); err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (d *SqliteBackend) GetBlockLossSummary(identifier string) (types.BlockLossSummary, error) {
	var summary types.BlockLossSummary
	err := d.db.QueryRow(`
				SELECT
				  block_height,
				  block_hash,
				  COUNT(DISTINCT(tx_id)) AS loss_tx_count,
				  SUM(amount) AS sum_of_losses
				FROM
				  losses
				WHERE
				  block_hash = ? OR block_height = ?
				GROUP BY
				  block_height, block_hash`, identifier, identifier).Scan(&summary.BlockHeight, &summary.BlockHash, &summary.LossOutputs, &summary.TotalLost)

	return summary, err
}

func (d *SqliteBackend) GetLatestBlock() (types.Block, error) {
	var block types.Block
	err := d.db.QueryRow("SELECT * FROM blocks ORDER BY block_height DESC LIMIT 1").Scan(&block.ID, &block.BlockHeight, &block.BlockHash, &block.BlockTimestamp, &block.ParentBlockHash, &block.NumTransactions, &block.CreatedAt)
	if err != nil {
		return types.Block{}, err
	}

	return block, nil
}

func (d *SqliteBackend) GetBlockIdentifiers(identifier string) (string, int64, error) {
	var blockHash string
	var blockHeight int64
	err := d.db.QueryRow("SELECT block_hash, block_height FROM blocks WHERE block_hash = ? OR block_height = ?", identifier, identifier).Scan(&blockHash, &blockHeight)

	return blockHash, blockHeight, err
}

func scanBlocks(rows *sql.Rows) ([]types.Block, error) {
	var blocks []types.Block
	for rows.Next() {
		var block types.Block
		err := rows.Scan(&block.ID, &block.BlockHeight, &block.BlockHash, &block.BlockTimestamp, &block.ParentBlockHash, &block.NumTransactions, &block.CreatedAt)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}
