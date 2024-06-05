package sqlite

import (
	"encoding/json"

	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetTransactionDetail(hash string) (types.TransactionDetail, error) {
	var transaction types.Transaction
	err := d.db.QueryRow("SELECT * FROM transactions WHERE tx_id = ?", hash).Scan(&transaction.TxID, &transaction.TransactionDetails, &transaction.BlockHash, &transaction.BlockHeight)
	if err != nil {
		return types.TransactionDetail{}, err
	}

	var txDetails types.TransactionDetail
	err = json.Unmarshal([]byte(transaction.TransactionDetails), &txDetails)
	if err != nil {
		return types.TransactionDetail{}, err
	}

	return txDetails, nil
}
