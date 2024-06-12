package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetTransaction(hash string) (types.Transaction, error) {
	var transaction types.Transaction
	err := d.db.QueryRow("SELECT * FROM transactions WHERE tx_id = ?", hash).Scan(&transaction.TxID, &transaction.TransactionDetails, &transaction.BlockHeight, &transaction.BlockHash)
	if err != nil {
		return types.Transaction{}, err
	}

	return transaction, nil
}

func (d *SqliteBackend) GetTransactionTxids(limit int, offset int) ([]string, error) {
	query := `SELECT tx_id FROM transaction_summary ORDER BY block_height ASC LIMIT ? OFFSET ?`

	rows, err := d.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var txids []string
	for rows.Next() {
		var txid string
		err := rows.Scan(&txid)
		if err != nil {
			return nil, err
		}

		txids = append(txids, txid)
	}

	return txids, nil
}

func (d *SqliteBackend) GetTransactionCount() (int, error) {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM transaction_summary").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (d *SqliteBackend) GetTransactionLossSummary(limit int) ([]types.TransactionLossSummary, error) {
	query := `SELECT tx_id, min(vout), sum(amount), block_height, block_hash FROM losses GROUP BY tx_id ORDER BY block_height DESC LIMIT ?`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	summaries, err := scanTransactionLossSummaries(rows)
	return summaries, err
}

func (d *SqliteBackend) GetTransactionLossSummaryForBlock(identifier string) ([]types.TransactionLossSummary, error) {
	query := `SELECT tx_id, min(vout), sum(amount), block_height, block_hash FROM losses WHERE block_hash = ? OR block_height = ? GROUP BY tx_id`

	rows, err := d.db.Query(query, identifier, identifier)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	summaries, err := scanTransactionLossSummaries(rows)
	return summaries, err
}

func (d *SqliteBackend) GetTransactionLossSummaryForScript(script string) ([]types.TransactionLossSummary, error) {
	query := `SELECT tx_id, min(vout), sum(amount), block_height, block_hash FROM losses WHERE burn_script = ? GROUP BY tx_id ORDER BY block_height DESC`

	rows, err := d.db.Query(query, script)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	summaries, err := scanTransactionLossSummaries(rows)
	return summaries, err
}

func (d *SqliteBackend) GetTransactionLossSummaryForTxid(txid string) (types.TransactionLossSummary, error) {
	query := `SELECT tx_id, min(vout), sum(amount), block_height, block_hash FROM losses WHERE tx_id = ?`

	var summary types.TransactionLossSummary
	err := d.db.QueryRow(query, txid).Scan(&summary.Txid, &summary.Vout, &summary.TotalLoss, &summary.BlockHeight, &summary.BlockHash)
	if err != nil {
		return types.TransactionLossSummary{}, err
	}

	return summary, nil
}

func scanTransactionLossSummaries(rows *sql.Rows) ([]types.TransactionLossSummary, error) {
	var summaries []types.TransactionLossSummary
	for rows.Next() {
		var summary types.TransactionLossSummary
		err := rows.Scan(&summary.Txid, &summary.Vout, &summary.TotalLoss, &summary.BlockHeight, &summary.BlockHash)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}
