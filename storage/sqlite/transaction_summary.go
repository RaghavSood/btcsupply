package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetTransactionSummary(limit int) ([]types.TransactionSummary, error) {
	query := `SELECT tx_id, coinbase, total_loss, block_height, block_hash FROM transaction_summary ORDER BY block_height DESC LIMIT ?`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	summaries, err := scanTransactionSummaries(rows)

	return summaries, err
}

func scanTransactionSummaries(rows *sql.Rows) ([]types.TransactionSummary, error) {
	var summaries []types.TransactionSummary
	for rows.Next() {
		var summary types.TransactionSummary
		err := rows.Scan(&summary.Txid, &summary.Coinbase, &summary.TotalLoss, &summary.BlockHeight, &summary.BlockHash)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}