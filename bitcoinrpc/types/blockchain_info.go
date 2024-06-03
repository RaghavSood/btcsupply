package types

type BlockchainInfo struct {
	Chain                string  `json:"chain"`
	Blocks               int64   `json:"blocks"`
	Headers              int64   `json:"headers"`
	Bestblockhash        string  `json:"bestblockhash"`
	Difficulty           float64 `json:"difficulty"`
	Time                 int     `json:"time"`
	Mediantime           int     `json:"mediantime"`
	Verificationprogress float64 `json:"verificationprogress"`
	Initialblockdownload bool    `json:"initialblockdownload"`
	Chainwork            string  `json:"chainwork"`
	SizeOnDisk           int64   `json:"size_on_disk"`
	Pruned               bool    `json:"pruned"`
	Warnings             string  `json:"warnings"`
}
