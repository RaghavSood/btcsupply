package types

type TxOutSetInfo struct {
	Height                 int64     `json:"height"`
	Bestblock              string    `json:"bestblock"`
	Txouts                 int       `json:"txouts"`
	Bogosize               int64     `json:"bogosize"`
	Muhash                 string    `json:"muhash"`
	TotalAmount            BTCString `json:"total_amount"`
	TotalUnspendableAmount BTCString `json:"total_unspendable_amount"`
	BlockInfo              BlockInfo `json:"block_info"`
}

type Unspendables struct {
	GenesisBlock     BTCString `json:"genesis_block"`
	Bip30            BTCString `json:"bip30"`
	Scripts          BTCString `json:"scripts"`
	UnclaimedRewards BTCString `json:"unclaimed_rewards"`
}

type BlockInfo struct {
	PrevoutSpent         BTCString    `json:"prevout_spent"`
	Coinbase             BTCString    `json:"coinbase"`
	NewOutputsExCoinbase BTCString    `json:"new_outputs_ex_coinbase"`
	Unspendable          BTCString    `json:"unspendable"`
	Unspendables         Unspendables `json:"unspendables"`
}
