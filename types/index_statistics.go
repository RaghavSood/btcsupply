package types

import "time"

type IndexStatistics struct {
	PlannedSupply     *BigInt   `json:"planned_supply"`
	CirculatingSupply *BigInt   `json:"circulating_supply"`
	BurnedSupply      *BigInt   `json:"burned_supply"`
	LastBlockHeight   int64     `json:"last_block_height"`
	LastBlockTime     time.Time `json:"last_block_time"`
	BurnOutputCount   int64     `json:"burned_output_count"`
	BurnScriptsCount  int64     `json:"burned_scripts_count"`
	CurrentPrice      float64   `json:"current_price"`
	AdjustedPrice     float64   `json:"adjusted_price"`
	PriceChange       float64   `json:"price_change"`
}
