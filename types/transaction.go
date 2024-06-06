package types

import "encoding/json"

type Transaction struct {
	TxID               string
	TransactionDetails string
	BlockHeight        int64
	BlockHash          string
}

// We store TransactionDetails in sqlite as a raw JSON to avoid having to handle
// various edge cases with variance in the transaction information, while still
// being able to support showing the full raw RPC response in the UI.
//
// This method will return the TransactionDetail struct from the raw JSON.
func (t Transaction) TransactionDetail() (TransactionDetail, error) {
	var txDetails TransactionDetail
	err := json.Unmarshal([]byte(t.TransactionDetails), &txDetails)
	if err != nil {
		return TransactionDetail{}, err
	}

	return txDetails, nil
}
