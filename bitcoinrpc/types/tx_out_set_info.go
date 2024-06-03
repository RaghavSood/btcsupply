package types

type TxOutSetInfo struct {
	Height                 int64     `json:"height"`
	Bestblock              string    `json:"bestblock"`
	Txouts                 int       `json:"txouts"`
	Bogosize               int64     `json:"bogosize"`
	Muhash                 string    `json:"muhash"`
	TotalAmount            float64   `json:"total_amount"`
	TotalUnspendableAmount float64   `json:"total_unspendable_amount"`
	BlockInfo              BlockInfo `json:"block_info"`
}

type Unspendables struct {
	GenesisBlock     float64 `json:"genesis_block"`
	Bip30            float64 `json:"bip30"`
	Scripts          float64 `json:"scripts"`
	UnclaimedRewards float64 `json:"unclaimed_rewards"`
}

type BlockInfo struct {
	PrevoutSpent         float64      `json:"prevout_spent"`
	Coinbase             float64      `json:"coinbase"`
	NewOutputsExCoinbase float64      `json:"new_outputs_ex_coinbase"`
	Unspendable          float64      `json:"unspendable"`
	Unspendables         Unspendables `json:"unspendables"`
}
