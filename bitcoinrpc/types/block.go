package types

type Block struct {
	Hash              string              `json:"hash"`
	Confirmations     int                 `json:"confirmations"`
	Height            int64               `json:"height"`
	Version           int                 `json:"version"`
	VersionHex        string              `json:"versionHex"`
	Merkleroot        string              `json:"merkleroot"`
	Time              int                 `json:"time"`
	Mediantime        int                 `json:"mediantime"`
	Nonce             int                 `json:"nonce"`
	Bits              string              `json:"bits"`
	Difficulty        float64             `json:"difficulty"`
	Chainwork         string              `json:"chainwork"`
	NTx               int                 `json:"nTx"`
	Previousblockhash string              `json:"previousblockhash"`
	Nextblockhash     string              `json:"nextblockhash"`
	Strippedsize      int                 `json:"strippedsize"`
	Size              int                 `json:"size"`
	Weight            int                 `json:"weight"`
	Tx                []TransactionDetail `json:"tx"`
}
