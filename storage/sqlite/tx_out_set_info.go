package sqlite

import (
	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetTxOutSetInfo(identifier string) (types.TxOutSetInfo, error) {
	var info types.TxOutSetInfo
	err := d.db.QueryRow("SELECT id, height, bestblock, txouts, bogosize, muhash, total_amount, total_unspendable_amount, prevout_spent, coinbase, new_outputs_ex_coinbase, unspendable, genesis_block, bip30, scripts, unclaimed_rewards, created_at FROM tx_out_set_info WHERE height = ? OR bestblock = ?", identifier, identifier).Scan(&info.ID, &info.Height, &info.Bestblock, &info.Txouts, &info.Bogosize, &info.Muhash, &info.TotalAmount, &info.TotalUnspendableAmount, &info.PrevoutSpent, &info.Coinbase, &info.NewOutputsExCoinbase, &info.Unspendable, &info.GenesisBlock, &info.Bip30, &info.Scripts, &info.UnclaimedRewards, &info.CreatedAt)
	if err != nil {
		return types.TxOutSetInfo{}, err
	}

	return info, nil
}
