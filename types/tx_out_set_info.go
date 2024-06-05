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
		TotalAmount:            FromBTCFloat64(info.TotalAmount),
		TotalUnspendableAmount: FromBTCFloat64(info.TotalUnspendableAmount),
		GenesisBlock:           FromBTCFloat64(info.BlockInfo.Unspendables.GenesisBlock),
		Bip30:                  FromBTCFloat64(info.BlockInfo.Unspendables.Bip30),
		Scripts:                FromBTCFloat64(info.BlockInfo.Unspendables.Scripts),
		UnclaimedRewards:       FromBTCFloat64(info.BlockInfo.Unspendables.UnclaimedRewards),
		PrevoutSpent:           FromBTCFloat64(info.BlockInfo.PrevoutSpent),
		Coinbase:               FromBTCFloat64(info.BlockInfo.Coinbase),
		NewOutputsExCoinbase:   FromBTCFloat64(info.BlockInfo.NewOutputsExCoinbase),
		Unspendable:            FromBTCFloat64(info.BlockInfo.Unspendable),
	}
}
