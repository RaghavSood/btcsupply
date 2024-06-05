package sqlite

import (
	btypes "github.com/RaghavSood/btcsupply/bitcoinrpc/types"
)

func (d *SqliteBackend) GetBlockStats(identifier string) (btypes.BlockStats, error) {
	var stats btypes.BlockStats
	err := d.db.QueryRow("SELECT * FROM block_stats WHERE block_hash = ? OR block_height = ?", identifier, identifier).Scan(&stats.Avgfee, &stats.Avgfeerate, &stats.Avgtxsize, &stats.Blockhash, &stats.Height, &stats.Ins, &stats.Maxfee, &stats.Maxfeerate, &stats.Maxtxsize, &stats.Medianfee, &stats.Mediantime, &stats.Mediantxsize, &stats.Minfee, &stats.Minfeerate, &stats.Mintxsize, &stats.Outs, &stats.Subsidy, &stats.SwtotalSize, &stats.SwtotalWeight, &stats.Swtxs, &stats.Time, &stats.TotalOut, &stats.TotalSize, &stats.TotalWeight, &stats.Totalfee, &stats.Txs, &stats.UtxoIncrease, &stats.UtxoSizeInc, &stats.UtxoIncreaseActual, &stats.UtxoSizeIncActual)
	if err != nil {
		return btypes.BlockStats{}, err
	}

	return stats, nil
}
