package storage

import (
	btypes "github.com/RaghavSood/btcsupply/bitcoinrpc/types"
	"github.com/RaghavSood/btcsupply/types"
)

type Storage interface {
	GetRecentLosses(limit int) ([]types.Loss, error)
	GetBlockLosses(hash string) ([]types.Loss, error)
	GetTransactionLosses(hash string) ([]types.Loss, error)

	GetLossyBlocks(limit int) ([]types.Block, error)
	GetBlock(hash string) (types.Block, error)
	GetLatestBlock() (types.Block, error)

	GetTransactionDetail(hash string) (types.TransactionDetail, error)

	RecordBlockIndexResults(block types.Block, txoutset types.TxOutSetInfo, blockstats btypes.BlockStats, losses []types.Loss, transactions []types.Transaction, spentTxids []string, spentVouts []int) error
	RecordTransactionIndexResults(losses []types.Loss, transactions []types.Transaction, spentTxids []string, spentVouts []int) error

	GetOnlyBurnScripts() ([]string, error)
	GetBurnScripts() ([]types.BurnScript, error)
	BurnScriptExists(script string) (bool, error)

	GetScriptQueue() ([]types.ScriptQueue, error)
	IncrementScriptQueueTryCount(scriptID int) error
	RecordScriptUnspents(script types.ScriptQueue, unspentTxids []string, unspentHeights []int64) error

	GetTransactionQueue() ([]types.TransactionQueue, error)
}
