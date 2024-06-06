package storage

import (
	btypes "github.com/RaghavSood/btcsupply/bitcoinrpc/types"
	"github.com/RaghavSood/btcsupply/types"
)

type Storage interface {
	GetRecentLosses(limit int) ([]types.Loss, error)
	GetBlockLosses(hash string) ([]types.Loss, error)
	GetTransactionLosses(hash string) ([]types.Loss, error)
	GetIndexStatistics(height int64) (types.IndexStatistics, error)

	GetTransactionLossSummary(limit int) ([]types.TransactionLossSummary, error)
	GetTransactionLossSummaryForBlock(identifier string) ([]types.TransactionLossSummary, error)

	GetLossyBlocks(limit int) ([]types.BlockLossSummary, error)
	GetBlock(identifier string) (types.Block, error)
	GetLatestBlock() (types.Block, error)
	GetBlockIdentifiers(identifier string) (string, int64, error)

	GetTxOutSetInfo(identifier string) (types.TxOutSetInfo, error)

	GetBlockStats(identifier string) (btypes.BlockStats, error)

	GetTransaction(hash string) (types.Transaction, error)

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
