package types

type ScriptGroupSummary struct {
	ScriptGroup  string
	TotalLoss    *BigInt
	Scripts      int
	Transactions int
}
