package types

import (
	"time"

	btypes "github.com/RaghavSood/btcsupply/bitcoinrpc/types"
)

type TxOutSetInfo struct {
	ID                     int       `json:"id"`
	Height                 int64     `json:"height"`
	Bestblock              string    `json:"bestblock"`
	Txouts                 int       `json:"txouts"`
	Bogosize               int64     `json:"bogosize"`
	Muhash                 string    `json:"muhash"`
	TotalAmount            *BigInt   `json:"total_amount"`
	TotalUnspendableAmount *BigInt   `json:"total_unspendable_amount"`
	GenesisBlock           *BigInt   `json:"genesis_block"`
	Bip30                  *BigInt   `json:"bip30"`
	Scripts                *BigInt   `json:"scripts"`
	UnclaimedRewards       *BigInt   `json:"unclaimed_rewards"`
	PrevoutSpent           *BigInt   `json:"prevout_spent"`
	Coinbase               *BigInt   `json:"coinbase"`
	NewOutputsExCoinbase   *BigInt   `json:"new_outputs_ex_coinbase"`
	Unspendable            *BigInt   `json:"unspendable"`
	CreatedAt              time.Time `json:"created_at"`
}

func FromRPCTxOutSetInfo(info btypes.TxOutSetInfo) TxOutSetInfo {
	return TxOutSetInfo{
		Height:                 info.Height,
		Bestblock:              info.Bestblock,
		Txouts:                 info.Txouts,
		Bogosize:               info.Bogosize,
		Muhash:                 info.Muhash,
		TotalAmount:            FromBTCString(BTCString(info.TotalAmount)),
		TotalUnspendableAmount: FromBTCString(BTCString(info.TotalUnspendableAmount)),
		GenesisBlock:           FromBTCString(BTCString(info.BlockInfo.Unspendables.GenesisBlock)),
		Bip30:                  FromBTCString(BTCString(info.BlockInfo.Unspendables.Bip30)),
		Scripts:                FromBTCString(BTCString(info.BlockInfo.Unspendables.Scripts)),
		UnclaimedRewards:       FromBTCString(BTCString(info.BlockInfo.Unspendables.UnclaimedRewards)),
		PrevoutSpent:           FromBTCString(BTCString(info.BlockInfo.PrevoutSpent)),
		Coinbase:               FromBTCString(BTCString(info.BlockInfo.Coinbase)),
		NewOutputsExCoinbase:   FromBTCString(BTCString(info.BlockInfo.NewOutputsExCoinbase)),
		Unspendable:            FromBTCString(BTCString(info.BlockInfo.Unspendable)),
	}
}
