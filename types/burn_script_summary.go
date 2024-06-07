package types

type BurnScriptSummary struct {
	Script          string
	ConfidenceLevel string
	Provenance      string
	Group           string
	TotalLoss       *BigInt
}
