package sqlite

import (
	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetTxOutSetInfo(identifier string) (types.TxOutSetInfo, error) {
	var info types.TxOutSetInfo
	err := d.db.QueryRow("SELECT * FROM txoutset_info WHERE height = ? OR bestblock = ?", identifier, identifier).Scan(&info.Height, &info.Bestblock, &info.Txouts, &info.Bogosize, &info.Muhash, &info.TotalAmount, &info.TotalUnspendableAmount, &info.GenesisBlock, &info.Bip30, &info.Scripts, &info.UnclaimedRewards, &info.PrevoutSpent, &info.Coinbase, &info.NewOutputsExCoinbase, &info.Unspendable)
	if err != nil {
		return types.TxOutSetInfo{}, err
	}

	return info, nil
}
