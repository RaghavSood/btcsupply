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
	GetTransactionLossSummaryForScript(script string) ([]types.TransactionLossSummary, error)
	GetTransactionLossSummaryForTxid(txid string) (types.TransactionLossSummary, error)

	GetTransactionSummary(limit int) ([]types.TransactionSummary, error)
	GetTransactionSummaryForTxid(txid string) (types.TransactionSummary, error)
	GetTransactionSummaryForBlock(identifier string) ([]types.TransactionSummary, error)

	GetLossyBlocks(limit int) ([]types.BlockLossSummary, error)
	GetBlockLossSummary(identifier string) (types.BlockLossSummary, error)
	GetBlock(identifier string) (types.Block, error)
	GetBlocksByHeights(heights []int64) ([]types.Block, error)
	GetLatestBlock() (types.Block, error)
	GetBlockIdentifiers(identifier string) (string, int64, error)

	GetTxOutSetInfo(identifier string) (types.TxOutSetInfo, error)

	GetBlockStats(identifier string) (btypes.BlockStats, error)

	GetTransaction(hash string) (types.Transaction, error)
	GetTransactionTxids(limit int, offset int) ([]string, error)
	GetTransactionCount() (int, error)

	RecordBlockIndexResults(block types.Block, txoutset types.TxOutSetInfo, blockstats btypes.BlockStats, losses []types.Loss, transactions []types.Transaction, spentTxids []string, spentVouts []int) error
	RecordTransactionIndexResults(losses []types.Loss, transactions []types.Transaction, spentTxids []string, spentVouts []int) error

	GetOnlyBurnScripts() ([]string, error)
	GetBurnScripts() ([]types.BurnScript, error)
	GetBurnScript(script string) (types.BurnScript, error)
	BurnScriptExists(script string) (bool, error)
	GetBurnScriptsByScripts(scripts []string) ([]types.BurnScript, error)
	GetUndecodedBurnScripts() ([]types.BurnScript, error)
	RecordDecodedBurnScript(script string, decodeScript string) error
	GetBurnScriptCount() (int, error)
	GetBurnScriptPage(limit int, offset int) ([]types.BurnScript, error)

	GetBurnScriptSummaries(limit int) ([]types.BurnScriptSummary, error)
	GetBurnScriptSummariesForGroup(group string) ([]types.BurnScriptSummary, error)
	GetBurnScriptSummary(script string) (types.BurnScriptSummary, error)

	GetScriptGroupSummaries(limit int) ([]types.ScriptGroupSummary, error)
	GetScriptGroupSummary(group string) (types.ScriptGroupSummary, error)

	GetOpReturnSummaries(limit int) ([]types.OpReturnSummary, error)
	GetOpReturnSummary(script string) (types.OpReturnSummary, error)

	GetScriptQueue() ([]types.ScriptQueue, error)
	IncrementScriptQueueTryCount(scriptID int) error
	RecordScriptUnspents(script types.ScriptQueue, unspentTxids []string, unspentHeights []int64) error

	GetTransactionQueue() ([]types.TransactionQueue, error)

	GetHeightLossSummary() ([]types.HeightLossSummary, error)
}
