package sqlite

import (
	btypes "github.com/RaghavSood/btcsupply/bitcoinrpc/types"
)

func (d *SqliteBackend) GetBlockStats(identifier string) (btypes.BlockStats, error) {
	var stats btypes.BlockStats
	err := d.db.QueryRow("SELECT avgfee, avgfeerate, avgtxsize, blockhash, height, ins, maxfee, maxfeerate, maxtxsize, medianfee, mediantime, mediantxsize, minfee, minfeerate, mintxsize, outs, subsidy, swtotal_size, swtotal_weight, swtxs, time, total_out, total_size, total_weight, totalfee, txs, utxo_increase, utxo_size_inc, utxo_increase_actual, utxo_size_inc_actual FROM block_stats WHERE blockhash = ? OR height = ?", identifier, identifier).Scan(&stats.Avgfee, &stats.Avgfeerate, &stats.Avgtxsize, &stats.Blockhash, &stats.Height, &stats.Ins, &stats.Maxfee, &stats.Maxfeerate, &stats.Maxtxsize, &stats.Medianfee, &stats.Mediantime, &stats.Mediantxsize, &stats.Minfee, &stats.Minfeerate, &stats.Mintxsize, &stats.Outs, &stats.Subsidy, &stats.SwtotalSize, &stats.SwtotalWeight, &stats.Swtxs, &stats.Time, &stats.TotalOut, &stats.TotalSize, &stats.TotalWeight, &stats.Totalfee, &stats.Txs, &stats.UtxoIncrease, &stats.UtxoSizeInc, &stats.UtxoIncreaseActual, &stats.UtxoSizeIncActual)
	if err != nil {
		return btypes.BlockStats{}, err
	}

	return stats, nil
}
