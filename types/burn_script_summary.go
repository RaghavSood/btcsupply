package types

import "encoding/json"

type BurnScriptSummary struct {
	Script          string
	ConfidenceLevel string
	Provenance      string
	Group           string
	DecodeScript    string
	Transactions    int
	TotalLoss       *BigInt
}

func (bss BurnScriptSummary) ParseDecodeScript() (DecodeScript, error) {
	if bss.DecodeScript == "" {
		return DecodeScript{
			Asm:     "",
			Type:    "nonstandard",
			Address: bss.Script,
			Script:  bss.Script,
		}, nil
	}

	var decodeScript DecodeScript
	err := json.Unmarshal([]byte(bss.DecodeScript), &decodeScript)

	return decodeScript, err
}
