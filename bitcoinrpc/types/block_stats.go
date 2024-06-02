package types

type BlockStats struct {
	Avgfee     int64  `json:"avgfee"`
	Avgfeerate int    `json:"avgfeerate"`
	Avgtxsize  int    `json:"avgtxsize"`
	Blockhash  string `json:"blockhash"`
	// We don't need this and it is too much effort to deal with this and sqlite
	// FeeratePercentiles []int  `json:"feerate_percentiles"`
	Height             int64 `json:"height"`
	Ins                int   `json:"ins"`
	Maxfee             int64 `json:"maxfee"`
	Maxfeerate         int   `json:"maxfeerate"`
	Maxtxsize          int   `json:"maxtxsize"`
	Medianfee          int64 `json:"medianfee"`
	Mediantime         int   `json:"mediantime"`
	Mediantxsize       int   `json:"mediantxsize"`
	Minfee             int64 `json:"minfee"`
	Minfeerate         int   `json:"minfeerate"`
	Mintxsize          int   `json:"mintxsize"`
	Outs               int   `json:"outs"`
	Subsidy            int64 `json:"subsidy"`
	SwtotalSize        int   `json:"swtotal_size"`
	SwtotalWeight      int   `json:"swtotal_weight"`
	Swtxs              int   `json:"swtxs"`
	Time               int   `json:"time"`
	TotalOut           int64 `json:"total_out"`
	TotalSize          int   `json:"total_size"`
	TotalWeight        int   `json:"total_weight"`
	Totalfee           int64 `json:"totalfee"`
	Txs                int   `json:"txs"`
	UtxoIncrease       int   `json:"utxo_increase"`
	UtxoSizeInc        int   `json:"utxo_size_inc"`
	UtxoIncreaseActual int   `json:"utxo_increase_actual"`
	UtxoSizeIncActual  int   `json:"utxo_size_inc_actual"`
}
