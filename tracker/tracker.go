package tracker

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/RaghavSood/btcsupply/bitcoinrpc"
	btypes "github.com/RaghavSood/btcsupply/bitcoinrpc/types"
	"github.com/RaghavSood/btcsupply/bloomfilter"
	"github.com/RaghavSood/btcsupply/btclogger"
	"github.com/RaghavSood/btcsupply/electrum"
	"github.com/RaghavSood/btcsupply/storage"
	"github.com/RaghavSood/btcsupply/types"
	"github.com/RaghavSood/btcsupply/util"
)

var log = btclogger.NewLogger("tracker")

type Tracker struct {
	db      storage.Storage
	client  *bitcoinrpc.RpcClient
	bf      *bloomfilter.BloomFilter
	eclient *electrum.Electrum
}

func NewTracker(db storage.Storage) *Tracker {
	eclient, err := electrum.NewElectrum()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to electrum server")
	}

	return &Tracker{
		db:      db,
		client:  bitcoinrpc.NewRpcClient(os.Getenv("BITCOIND_HOST"), os.Getenv("BITCOIND_USER"), os.Getenv("BITCOIND_PASS")),
		bf:      bloomfilter.NewBloomFilter(),
		eclient: eclient,
	}
}

func (t *Tracker) Run() {
	log.Info().Msg("Starting tracker")

	log.Info().Msg("Loading burn scripts from database")
	burnScripts, err := t.db.GetOnlyBurnScripts()
	if err == sql.ErrNoRows {
		log.Info().Msg("No burn scripts found in database")
		err = nil
	}
	if err != nil {
		log.Error().Err(err).Msg("Failed to load burn scripts")
	} else {
		log.Info().Int("count", len(burnScripts)).Msg("Loading burn scripts into bloom filter")
		t.bf.AddStrings(burnScripts)
		log.Info().Int("count", len(burnScripts)).Msg("Burn scripts loaded")
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Info().Msg("Checking for new blocks")
			info, err := t.client.GetBlockchainInfo()
			if err != nil {
				log.Error().Err(err).Msg("Failed to get blockchain info")
				continue
			}

			latestBlock, err := t.db.GetLatestBlock()
			if err == sql.ErrNoRows {
				latestBlock = types.Block{
					BlockHeight: -1,
				}
				log.Info().Msg("No blocks found in database, starting initial sync")
				err = nil
			}
			if err != nil {
				log.Error().Err(err).Msg("Failed to get latest block")
				continue
			}

			log.Info().
				Int64("latest_block", latestBlock.BlockHeight).
				Int64("current_block", info.Blocks).
				Msg("Checking for new blocks")

			t.processScriptQueue()
			t.processTransactionQueue()

			// We limit ourselves to batch processing 10 blocks at a time
			// so that other indexing jobs also run often enough
			target := min(latestBlock.BlockHeight+1+10, info.Blocks)

			for i := latestBlock.BlockHeight + 1; i <= target; i++ {
				err = t.processBlock(i)
				if err != nil {
					log.Error().Err(err).Int64("block_height", i).Msg("Failed to process block")
					break
				}
			}
		}
	}
}

func (t *Tracker) processTransactionQueue() {
	log.Info().Msg("Processing transaction queue")

	txs, err := t.db.GetTransactionQueue()
	if err == sql.ErrNoRows {
		log.Info().Msg("No transactions in queue")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("Failed to get transactions from queue")
		return
	}

	for _, tx := range txs {
		log.Info().Str("txid", tx.Txid).Msg("Processing transaction")
		txDetails, err := t.client.GetTransaction(tx.Txid)
		if err != nil {
			log.Error().Err(err).Str("txid", tx.Txid).Msg("Failed to get transaction details")
			continue
		}

		if txDetails.Blockhash == "" {
			log.Info().Str("txid", tx.Txid).Msg("Transaction not yet confirmed")
			continue
		}

		losses, txs, spentTxids, spentVouts := t.scanTransactions(txDetails.Blockhash, tx.BlockHeight, []btypes.TransactionDetail{txDetails})
		log.Info().
			Str("txid", tx.Txid).
			Str("blockhash", txDetails.Blockhash).
			Int64("block_height", tx.BlockHeight).
			Int("losses", len(losses)).
			Int("transactions", len(txs)).
			Int("spent_txids", len(spentTxids)).
			Int("spent_vouts", len(spentVouts)).
			Msg("Transaction scan results")

		err = t.db.RecordTransactionIndexResults(losses, txs, spentTxids, spentVouts)
		if err != nil {
			log.Error().Err(err).Str("txid", tx.Txid).Msg("Failed to record transaction index results")
		}
	}
}

func (t *Tracker) processScriptQueue() {
	log.Info().Msg("Processing script queue")

	scripts, err := t.db.GetScriptQueue()
	if err == sql.ErrNoRows {
		log.Info().Msg("No scripts in queue")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("Failed to get addresses from queue")
		return
	}

	for _, script := range scripts {
		if script.TryCount > 5 {
			log.Warn().Str("script", script.Script).Int("try_count", script.TryCount).Msg("Script has been tried too many times, skipping")
			continue
		}

		log.Info().Str("script", script.Script).Int("try_count", script.TryCount).Msg("Processing script")
		unspentTxids, unspentHeights, err := t.eclient.GetScriptUnspents(script.Script)
		if err != nil {
			log.Error().Err(err).Str("script", script.Script).Msg("Failed to get unspents")
			err = t.db.IncrementScriptQueueTryCount(script.ID)
			if err != nil {
				log.Error().Err(err).Int("id", script.ID).Msg("Failed to increment try count")
			}
			continue
		}

		err = t.db.RecordScriptUnspents(script, unspentTxids, unspentHeights)
		if err != nil {
			log.Error().Err(err).Str("script", script.Script).Msg("Failed to record unspents")
			err = t.db.IncrementScriptQueueTryCount(script.ID)
			if err != nil {
				log.Error().Err(err).Int("id", script.ID).Msg("Failed to increment try count")
			}
		}
	}
}

func (t *Tracker) processBlock(height int64) error {
	log.Info().Int64("block_height", height).Msg("Processing block")
	blockStats, err := t.client.GetBlockStats(int64(height))
	if err != nil {
		return err
	}

	log.Info().
		Int64("subsidy", blockStats.Subsidy).
		Str("hash", blockStats.Blockhash).
		Int64("totalfee", blockStats.Totalfee).
		Int64("height", blockStats.Height).
		Msg("Block stats")

	coinStats, err := t.client.GetTxOutSetInfo("muhash", int64(height), true)
	if err != nil {
		return err
	}

	log.Info().
		Float64("total_amount", coinStats.TotalAmount).
		Float64("total_unspendable_amount", coinStats.TotalUnspendableAmount).
		Str("bestblock", coinStats.Bestblock).
		Float64("coinbase", coinStats.BlockInfo.Coinbase).
		Float64("unspendable", coinStats.BlockInfo.Unspendable).
		Float64("genesis_block", coinStats.BlockInfo.Unspendables.GenesisBlock).
		Float64("bip30", coinStats.BlockInfo.Unspendables.Bip30).
		Float64("scripts", coinStats.BlockInfo.Unspendables.Scripts).
		Float64("unclaimed_rewards", coinStats.BlockInfo.Unspendables.UnclaimedRewards).
		Int64("height", coinStats.Height).
		Msg("Coin stats")

	block, err := t.client.GetBlock(blockStats.Blockhash)
	if err != nil {
		return err
	}

	log.Info().
		Int("nTx", block.NTx).
		Int64("height", block.Height).
		Msg("Block")

	var losses []types.Loss
	var transactions []types.Transaction

	feesAccumulated := blockStats.Totalfee
	coinbaseEntitlement := blockStats.Subsidy
	expectedCoinbase := coinbaseEntitlement + feesAccumulated

	coinbaseMinted := util.FloatBTCToSats(coinStats.BlockInfo.Coinbase)

	if coinbaseMinted != expectedCoinbase {
		log.Warn().
			Int64("coinbase_minted", coinbaseMinted).
			Int64("coinbase_entitlement", coinbaseEntitlement).
			Msg("Coinbase mismatch")

		losses = append(losses, types.Loss{
			TxID:      block.Tx[0].Txid,
			BlockHash: block.Hash,
			Vout:      -1,
			Amount:    types.FromMathBigInt(big.NewInt(expectedCoinbase - coinbaseMinted)),
		})

		jsonTx, err := json.Marshal(block.Tx[0])
		if err != nil {
			return fmt.Errorf("failed to marshal coinbase tx: %v", err)
		}

		transactions = append(transactions, types.Transaction{
			TxID:               block.Tx[0].Txid,
			TransactionDetails: string(jsonTx),
		})
	}

	if coinStats.BlockInfo.Unspendables.Scripts != 0 {
		log.Warn().Float64("scripts", coinStats.BlockInfo.Unspendables.Scripts).Msg("Unspendable scripts")
	}

	// The Genesis block has a unique case where the funds are lost both because
	// they are sent to Satoshi's keys AND an underlying bug in the code
	//
	// We don't want to double count it, so we skip this step for the genesis block
	var txLosses []types.Loss
	var txTransactions []types.Transaction
	var spentTxids []string
	var spentVouts []int
	if block.Height > 0 {
		txLosses, txTransactions, spentTxids, spentVouts = t.scanTransactions(block.Hash, block.Height, block.Tx)

		losses = append(losses, txLosses...)
		transactions = append(transactions, txTransactions...)
	}

	err = t.db.RecordBlockIndexResults(types.FromRPCBlock(block), types.FromRPCTxOutSetInfo(coinStats), blockStats, losses, transactions, spentTxids, spentVouts)
	if err != nil {
		return fmt.Errorf("failed to record block index results: %v", err)
	}

	return nil
}

func (t *Tracker) scanTransactions(blockhash string, blockHeight int64, transactions []btypes.TransactionDetail) ([]types.Loss, []types.Transaction, []string, []int) {
	var losses []types.Loss
	var txs []types.Transaction
	var spentTxids []string
	var spentVouts []int

	for _, tx := range transactions {
		for _, vin := range tx.Vin {
			if vin.Coinbase != "" {
				continue
			}

			spentScript := vin.Prevout.ScriptPubKey.Hex

			if t.bf.TestString(spentScript) {
				exists, err := t.db.BurnScriptExists(spentScript)
				if err != nil {
					log.Error().Err(err).Str("script", spentScript).Msg("Failed to check if script exists")
					continue
				}

				log.Info().
					Str("script", spentScript).
					Bool("exists", exists).
					Str("txid", vin.Txid).
					Int("vout", vin.Vout).
					Msg("Identified spending of burn script output")

				if exists {
					spentTxids = append(spentTxids, vin.Txid)
					spentVouts = append(spentVouts, vin.Vout)
				}
			}
		}

		for _, vout := range tx.Vout {
			if t.bf.TestString(vout.ScriptPubKey.Hex) {
				exists, err := t.db.BurnScriptExists(vout.ScriptPubKey.Hex)
				if err != nil {
					log.Error().Err(err).Str("script", vout.ScriptPubKey.Hex).Msg("Failed to check if script exists")
					continue
				}

				log.Info().Str("script", vout.ScriptPubKey.Hex).Bool("exists", exists).Msg("Burn script identified")

				if exists {
					jsonTx, err := json.Marshal(tx)
					if err != nil {
						log.Error().Err(err).Str("txid", tx.Txid).Msg("Failed to marshal transaction")
						continue
					}

					losses = append(losses, types.Loss{
						TxID:        tx.Txid,
						BlockHash:   blockhash,
						BlockHeight: blockHeight,
						Vout:        vout.N,
						Amount:      types.FromBTCFloat64(vout.Value),
					})

					txs = append(txs, types.Transaction{
						TxID:               tx.Txid,
						TransactionDetails: string(jsonTx),
						BlockHeight:        blockHeight,
						BlockHash:          blockhash,
					})
				}
			}
		}
	}

	return losses, txs, spentTxids, spentVouts
}
