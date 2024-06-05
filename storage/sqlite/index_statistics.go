package sqlite

import (
	"fmt"

	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetIndexStatistics(height int64) (types.IndexStatistics, error) {
	var stats types.IndexStatistics
	var lastBlock types.Block
	var err error
	if height < 0 {
		lastBlock, err = d.GetLatestBlock()
		if err != nil {
			return stats, err
		}
	} else {
		lastBlock, err = d.GetBlock(fmt.Sprintf("%d", height))
		if err != nil {
			return stats, err
		}
	}

	stats.LastBlockHeight = lastBlock.BlockHeight
	stats.LastBlockTime = lastBlock.BlockTimestamp

	err = d.db.QueryRow("SELECT COUNT(*), SUM(amount) FROM losses WHERE block_height <= ?", lastBlock.BlockHeight).Scan(&stats.BurnOutputCount, &stats.BurnedSupply)
	if err != nil {
		return stats, err
	}

	err = d.db.QueryRow("SELECT COUNT(*) FROM burn_scripts").Scan(&stats.BurnScriptsCount)
	if err != nil {
		return stats, err
	}

	return stats, nil
}
