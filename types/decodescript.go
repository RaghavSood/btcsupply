package types

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
