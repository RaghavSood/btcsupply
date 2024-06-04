package sqlite

import (
	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetIndexStatistics() (types.IndexStatistics, error) {
	var stats types.IndexStatistics
	lastBlock, err := d.GetLatestBlock()
	if err != nil {
		return stats, err
	}

	stats.LastBlockHeight = lastBlock.BlockHeight
	stats.LastBlockTime = lastBlock.BlockTimestamp

	err = d.db.QueryRow("SELECT COUNT(*) FROM losses").Scan(&stats.BurnOutputCount)
	if err != nil {
		return stats, err
	}

	err = d.db.QueryRow("SELECT COUNT(*) FROM burn_scripts").Scan(&stats.BurnScriptsCount)
	if err != nil {
		return stats, err
	}

	err = d.db.QueryRow("SELECT SUM(amount) FROM losses").Scan(&stats.BurnedSupply)
	if err != nil {
		return stats, err
	}

	return stats, nil
}
