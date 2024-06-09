package types

import (
	"encoding/json"
	"time"
)

type BurnScript struct {
	ID              int       `json:"id"`
	Script          string    `json:"script"`
	ConfidenceLevel string    `json:"confidence_level"`
	Provenance      string    `json:"provenance"`
	CreatedAt       time.Time `json:"created_at"`
	ScriptGroup     string    `json:"script_group"`
	DecodeScript    string    `json:"decodescript"`
}

func (bs BurnScript) ParseDecodeScript() (DecodeScript, error) {
	if bs.DecodeScript == "" {
		return DecodeScript{
			Asm:     "",
			Type:    "nonstandard",
			Address: bs.Script,
			Script:  bs.Script,
		}, nil
	}

	var decodeScript DecodeScript
	err := json.Unmarshal([]byte(bs.DecodeScript), &decodeScript)

	return decodeScript, err
}
