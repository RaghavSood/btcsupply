package sqlite

import (
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

func (d *SqliteBackend) GetTransactionLossSummary(limit int) ([]types.TransactionLossSummary, error) {
	query := `SELECT tx_id, sum(amount), block_height, block_hash FROM losses GROUP BY tx_id ORDER BY block_height DESC LIMIT ?`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var summaries []types.TransactionLossSummary
	for rows.Next() {
		var summary types.TransactionLossSummary
		err = rows.Scan(&summary.Txid, &summary.TotalLoss, &summary.BlockHeight, &summary.BlockHash)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}
