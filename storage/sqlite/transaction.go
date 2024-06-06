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
