package types

type OpReturnSummary struct {
	Script       string `json:"script"`
	Transactions int
	TotalLoss    *BigInt
}
