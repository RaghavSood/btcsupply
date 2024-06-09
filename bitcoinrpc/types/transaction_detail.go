package types

type TransactionDetail struct {
	Txid      string    `json:"txid"`
	Hash      string    `json:"hash"`
	Version   int       `json:"version"`
	Size      int       `json:"size"`
	Vsize     int       `json:"vsize"`
	Weight    int       `json:"weight"`
	Locktime  int       `json:"locktime"`
	Vin       []Vin     `json:"vin"`
	Vout      []Vout    `json:"vout"`
	Fee       BTCString `json:"fee"`
	Hex       string    `json:"hex"`
	Blockhash string    `json:"blockhash,omitempty"`
	Time      int       `json:"time,omitempty"`
	Blocktime int       `json:"blocktime,omitempty"`
}

type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type Prevout struct {
	Generated    bool         `json:"generated"`
	Height       int          `json:"height"`
	Value        BTCString    `json:"value"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

type ScriptPubKey struct {
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address,omitempty"`
	Type    string `json:"type"`
}

type Vout struct {
	Value        BTCString    `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

type Vin struct {
	Txid        string    `json:"txid"`
	Vout        int       `json:"vout"`
	ScriptSig   ScriptSig `json:"scriptSig"`
	Txinwitness []string  `json:"txinwitness,omitempty"`
	Prevout     Prevout   `json:"prevout"`
	Coinbase    string    `json:"coinbase,omitempty"`
	Sequence    int64     `json:"sequence"`
}
