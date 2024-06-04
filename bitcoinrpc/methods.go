package bitcoinrpc

import (
	"encoding/json"
	"fmt"

	"github.com/RaghavSood/btcsupply/bitcoinrpc/types"
)

func (rpc *RpcClient) GetTxOutSetInfo(hashType string, blockHeight int64, useIndex bool) (types.TxOutSetInfo, error) {
	result, err := rpc.Do("gettxoutsetinfo", []interface{}{hashType, blockHeight, useIndex})
	if err != nil {
		return types.TxOutSetInfo{}, err
	}

	var stats types.TxOutSetInfo
	if err := json.Unmarshal(result, &stats); err != nil {
		return types.TxOutSetInfo{}, fmt.Errorf("failed to unmarshal gettxoutsetinfo response: %v", err)
	}

	return stats, nil
}

func (rpc *RpcClient) GetBlockchainInfo() (types.BlockchainInfo, error) {
	result, err := rpc.Do("getblockchaininfo", []interface{}{})
	if err != nil {
		return types.BlockchainInfo{}, err
	}

	var info types.BlockchainInfo
	if err := json.Unmarshal(result, &info); err != nil {
		return types.BlockchainInfo{}, fmt.Errorf("failed to unmarshal getblockchaininfo response: %v", err)
	}

	return info, nil
}

func (rpc *RpcClient) GetBlockStats(height int64) (types.BlockStats, error) {
	result, err := rpc.Do("getblockstats", []interface{}{height})
	if err != nil {
		return types.BlockStats{}, err
	}

	var stats types.BlockStats
	if err := json.Unmarshal(result, &stats); err != nil {
		return types.BlockStats{}, fmt.Errorf("failed to unmarshal getblockstats response: %v", err)
	}

	return stats, nil
}

func (rpc *RpcClient) GetBlock(hash string) (types.Block, error) {
	// Ensure we always get the block in the most verbose mode
	result, err := rpc.Do("getblock", []interface{}{hash, 3})
	if err != nil {
		return types.Block{}, err
	}

	var block types.Block
	if err := json.Unmarshal(result, &block); err != nil {
		return types.Block{}, fmt.Errorf("failed to unmarshal getblock response: %v", err)
	}

	return block, nil
}

func (rpc *RpcClient) GetTransaction(txid string) (types.TransactionDetail, error) {
	result, err := rpc.Do("getrawtransaction", []interface{}{txid, 2})
	if err != nil {
		return types.Transaction{}, err
	}

	var tx types.Transaction
	if err := json.Unmarshal(result, &tx); err != nil {
		return types.Transaction{}, fmt.Errorf("failed to unmarshal getrawtransaction response: %v", err)
	}

	return tx, nil
}
