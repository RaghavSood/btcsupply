package types

type BurnScriptSummary struct {
	Script          string
	ConfidenceLevel string
	Provenance      string
	Group           string
	DecodeScript    string
	Transactions    int
	TotalLoss       *BigInt
}
