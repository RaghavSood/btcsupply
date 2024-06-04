package sqlite

import (
	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetTransactionQueue() ([]types.TransactionQueue, error) {
	rows, err := d.db.Query("SELECT * FROM transaction_queue")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []types.TransactionQueue
	for rows.Next() {
		var transaction types.TransactionQueue
		err := rows.Scan(&transaction.ID, &transaction.Txid, &transaction.BlockHeight, &transaction.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
