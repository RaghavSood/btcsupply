package types

import (
	"github.com/RaghavSood/btcsupply/notes"
	"golang.org/x/exp/maps"
)

type TransactionDetail struct {
	Txid      string  `json:"txid"`
	Hash      string  `json:"hash"`
	Version   int     `json:"version"`
	Size      int     `json:"size"`
	Vsize     int     `json:"vsize"`
	Weight    int     `json:"weight"`
	Locktime  int     `json:"locktime"`
	Vin       []Vin   `json:"vin"`
	Vout      []Vout  `json:"vout"`
	Fee       float64 `json:"fee"`
	Hex       string  `json:"hex"`
	Blockhash string  `json:"blockhash,omitempty"`
	Time      int     `json:"time,omitempty"`
	Blocktime int     `json:"blocktime,omitempty"`
}

type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type Prevout struct {
	Generated    bool         `json:"generated"`
	Height       int          `json:"height"`
	Value        float64      `json:"value"`
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
	Value        float64      `json:"value"`
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

func (t *TransactionDetail) NotePointers() ([]notes.NotePointer, bool, []string) {
	hasNulldata := false
	scriptsSeen := make(map[string]bool)
	notePointers := make([]notes.NotePointer, 0)

	for _, vout := range t.Vout {
		noteType := notes.Script

		if vout.ScriptPubKey.Type == "nulldata" {
			hasNulldata = true
			noteType = notes.NullData
		}

		notePointers = append(notePointers, notes.NotePointer{
			NoteType:     noteType,
			PathElements: []string{vout.ScriptPubKey.Hex},
		})

		scriptsSeen[vout.ScriptPubKey.Hex] = true
	}

	scripts := maps.Keys(scriptsSeen)
	return notePointers, hasNulldata, scripts
}

func ScriptPubKeyDisplay(scriptPubKey ScriptPubKey) string {
	if scriptPubKey.Address != "" {
		return scriptPubKey.Address
	}
	return scriptPubKey.Hex
}

func ValueToBigInt(value float64) *BigInt {
	return FromBTCFloat64(value)
}
