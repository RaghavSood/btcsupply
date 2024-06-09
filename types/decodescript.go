package types

import "encoding/json"

type DecodeScript struct {
	Script  string `json:"script"`
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Type    string `json:"type"`
	Address string `json:"address"`
	P2SH    string `json:"p2sh"`
	Segwit  Segwit `json:"segwit"`
}

type Segwit struct {
	Asm        string `json:"asm"`
	Hex        string `json:"hex"`
	Type       string `json:"type"`
	Address    string `json:"address"`
	Desc       string `json:"desc"`
	P2SHSegwit string `json:"p2sh-segwit"`
}

func (ds DecodeScript) DisplayAddress(fallback string) string {
	if ds.Address != "" {
		return ds.Address
	}

	result := ds.Script
	if result == "" {
		result = fallback
	}

	return result
}

func ParseDecodeScriptJSON(jsonString string) (DecodeScript, error) {
	if jsonString == "" {
		return DecodeScript{}, nil
	}
	var ds DecodeScript
	err := json.Unmarshal([]byte(jsonString), &ds)
	if err != nil {
		return ds, err
	}

	return ds, nil
}
