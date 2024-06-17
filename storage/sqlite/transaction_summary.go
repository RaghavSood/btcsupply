package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetTransactionSummary(limit int, minLoss int64) ([]types.TransactionSummary, error) {
	query := `SELECT ts.tx_id, ts.coinbase, ts.total_loss, ts.block_height, ts.block_hash, b.block_timestamp FROM transaction_summary ts JOIN blocks b ON ts.block_hash = b.block_hash WHERE total_loss >= ? ORDER BY ts.block_height DESC LIMIT ?`

	rows, err := d.db.Query(query, minLoss, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	summaries, err := scanTransactionSummaries(rows)

	return summaries, err
}

func (d *SqliteBackend) GetTransactionSummaryForBlock(identifier string) ([]types.TransactionSummary, error) {
	query := `SELECT ts.tx_id, ts.coinbase, ts.total_loss, ts.block_height, ts.block_hash, b.block_timestamp FROM transaction_summary ts JOIN blocks b ON ts.block_hash = b.block_hash WHERE ts.block_hash = ? OR ts.block_height = ? ORDER BY ts.total_loss DESC`

	rows, err := d.db.Query(query, identifier, identifier)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	summaries, err := scanTransactionSummaries(rows)

	return summaries, err
}

func (d *SqliteBackend) GetTransactionSummaryForTxid(txid string) (types.TransactionSummary, error) {
	query := `SELECT ts.tx_id, ts.coinbase, ts.total_loss, ts.block_height, ts.block_hash, b.block_timestamp FROM transaction_summary ts JOIN blocks b ON ts.block_hash = b.block_hash WHERE ts.tx_id = ?`

	var summary types.TransactionSummary
	err := d.db.QueryRow(query, txid).Scan(&summary.Txid, &summary.Coinbase, &summary.TotalLoss, &summary.BlockHeight, &summary.BlockHash)
	if err != nil {
		return types.TransactionSummary{}, err
	}

	return summary, nil
}

func scanTransactionSummaries(rows *sql.Rows) ([]types.TransactionSummary, error) {
	var summaries []types.TransactionSummary
	for rows.Next() {
		var summary types.TransactionSummary
		err := rows.Scan(&summary.Txid, &summary.Coinbase, &summary.TotalLoss, &summary.BlockHeight, &summary.BlockHash, &summary.BlockTimestamp)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}
