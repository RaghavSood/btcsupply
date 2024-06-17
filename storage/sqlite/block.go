package sqlite

import (
	"database/sql"
	"strings"

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

func (d *SqliteBackend) GetBlocksByHeights(heights []int64) ([]types.Block, error) {
	anyHeights := make([]interface{}, len(heights))
	for i, height := range heights {
		anyHeights[i] = height
	}

	rows, err := d.db.Query("SELECT * FROM blocks WHERE block_height IN (?"+strings.Repeat(",?", len(heights)-1)+")", anyHeights...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanBlocks(rows)
}

func (d *SqliteBackend) GetLossyBlocks(limit int) ([]types.BlockLossSummary, error) {
	query := `SELECT ts.block_height, ts.block_hash, COUNT(*), SUM(ts.total_loss), b.block_timestamp FROM transaction_summary ts JOIN blocks b ON ts.block_hash = b.block_hash GROUP BY ts.block_height ORDER BY ts.block_height DESC LIMIT ?`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []types.BlockLossSummary
	for rows.Next() {
		var summary types.BlockLossSummary
		if err := rows.Scan(
			&summary.BlockHeight, &summary.BlockHash, &summary.LossOutputs, &summary.TotalLost, &summary.BlockTimestamp); err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (d *SqliteBackend) GetLossyBlocksWithMinimumLoss(limit int, minimumLoss int64) ([]types.BlockLossSummary, error) {
	query := `SELECT ts.block_height, ts.block_hash, COUNT(*), SUM(ts.total_loss) as block_loss, b.block_timestamp FROM transaction_summary ts JOIN blocks b ON ts.block_hash = b.block_hash GROUP BY ts.block_height HAVING block_loss >= ? ORDER BY ts.block_height DESC LIMIT ?`

	rows, err := d.db.Query(query, minimumLoss, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []types.BlockLossSummary
	for rows.Next() {
		var summary types.BlockLossSummary
		if err := rows.Scan(
			&summary.BlockHeight, &summary.BlockHash, &summary.LossOutputs, &summary.TotalLost, &summary.BlockTimestamp); err != nil {
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
				  l.block_height,
				  l.block_hash,
				  COUNT(DISTINCT(l.tx_id)) AS loss_tx_count,
				  SUM(l.amount) AS sum_of_losses,
					b.block_timestamp
				FROM
				  losses l
				JOIN
				  blocks b ON l.block_hash = b.block_hash
				WHERE
				 l.block_hash = ? OR l.block_height = ?
				GROUP BY
				  l.block_height`, identifier, identifier).Scan(&summary.BlockHeight, &summary.BlockHash, &summary.LossOutputs, &summary.TotalLost, &summary.BlockTimestamp)

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
