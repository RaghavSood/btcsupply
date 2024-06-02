package bitcoinrpc

import (
	"encoding/json"
	"fmt"

	"github.com/RaghavSood/btcsupply/bitcoinrpc/types"
)

func (rpc *RpcClient) GetTxOutSetInfo(hashType string, blockHeight int64, useIndex bool) (types.Coinstats, error) {
	result, err := rpc.Do("gettxoutsetinfo", []interface{}{hashType, blockHeight, useIndex})
	if err != nil {
		return types.Coinstats{}, err
	}

	var stats types.Coinstats
	if err := json.Unmarshal(result, &stats); err != nil {
		return types.Coinstats{}, fmt.Errorf("failed to unmarshal gettxoutsetinfo response: %v", err)
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
