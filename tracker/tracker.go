package tracker

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/RaghavSood/btcsupply/bitcoinrpc"
	btypes "github.com/RaghavSood/btcsupply/bitcoinrpc/types"
	"github.com/RaghavSood/btcsupply/bloomfilter"
	"github.com/RaghavSood/btcsupply/btclogger"
	"github.com/RaghavSood/btcsupply/electrum"
	"github.com/RaghavSood/btcsupply/storage"
	"github.com/RaghavSood/btcsupply/types"
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

	ticker := time.NewTicker(5 * time.Second)
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

			t.decodedBurnScripts()
			t.processScriptQueue()
			t.processTransactionQueue()

			// We limit ourselves to batch processing 50 blocks at a time
			// so that other indexing jobs also run often enough
			target := min(latestBlock.BlockHeight+1+50, info.Blocks)

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

func (t *Tracker) decodedBurnScripts() {
	undecodedBurnScripts, err := t.db.GetUndecodedBurnScripts()
	if err == sql.ErrNoRows {
		log.Info().Msg("No undecoded burn scripts found")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("Failed to get undecoded burn scripts")
		return
	}

	for _, script := range undecodedBurnScripts {
		decoded, err := t.client.DecodeScript(script.Script)
		if err != nil {
			log.Error().Err(err).Str("script", script.Script).Msg("Failed to decode script")
			continue
		}

		decoded.Script = script.Script

		jsonScript, err := json.Marshal(decoded)
		if err != nil {
			log.Error().
				Err(err).
				Str("script", script.Script).
				Msg("Failed to marshal decoded script")
			continue
		}

		log.Info().Str("script", script.Script).Msg("Decoded script")

		err = t.db.RecordDecodedBurnScript(script.Script, string(jsonScript))
		if err != nil {
			log.Error().Err(err).Str("script", script.Script).Msg("Failed to update script decode")
			continue
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
	start := time.Now()
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
		Stringer("total_amount", coinStats.TotalAmount).
		Stringer("total_unspendable_amount", coinStats.TotalUnspendableAmount).
		Str("bestblock", coinStats.Bestblock).
		Stringer("coinbase", coinStats.BlockInfo.Coinbase).
		Stringer("unspendable", coinStats.BlockInfo.Unspendable).
		Stringer("genesis_block", coinStats.BlockInfo.Unspendables.GenesisBlock).
		Stringer("bip30", coinStats.BlockInfo.Unspendables.Bip30).
		Stringer("scripts", coinStats.BlockInfo.Unspendables.Scripts).
		Stringer("unclaimed_rewards", coinStats.BlockInfo.Unspendables.UnclaimedRewards).
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

	coinbaseMinted := types.FromBTCString(types.BTCString(coinStats.BlockInfo.Coinbase))
	expectedCoinbaseBig := types.FromMathBigInt(big.NewInt(expectedCoinbase))

	if coinbaseMinted.Cmp(expectedCoinbaseBig.BigInt()) != 0 {
		log.Warn().
			Stringer("coinbase_minted", coinbaseMinted).
			Stringer("coinbase_entitlement", expectedCoinbaseBig).
			Msg("Coinbase mismatch")

		lostAmount := big.NewInt(0).Sub(expectedCoinbaseBig.BigInt(), coinbaseMinted.BigInt())

		losses = append(losses, types.Loss{
			TxID:        block.Tx[0].Txid,
			BlockHash:   block.Hash,
			BlockHeight: block.Height,
			Vout:        -1,
			Amount:      types.FromMathBigInt(lostAmount),
		})

		jsonTx, err := json.Marshal(block.Tx[0])
		if err != nil {
			return fmt.Errorf("failed to marshal coinbase tx: %v", err)
		}

		transactions = append(transactions, types.Transaction{
			TxID:               block.Tx[0].Txid,
			TransactionDetails: string(jsonTx),
			BlockHeight:        block.Height,
			BlockHash:          block.Hash,
		})
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
		start := time.Now()
		txLosses, txTransactions, spentTxids, spentVouts = t.scanTransactions(block.Hash, block.Height, block.Tx)
		elapsed := time.Since(start)
		fastStart := time.Now()
		fastTxLosses, fastTxTransactions, fastSpentTxids, fastSpentVouts := t.fastScanTransactions(block.Hash, block.Height, block.Tx)
		fastElapsed := time.Since(fastStart)

		log.Info().
			Stringer("slow_elapsed", elapsed).
			Stringer("fast_elapsed", fastElapsed).
			Int("slow_losses", len(txLosses)).
			Int("fast_losses", len(fastTxLosses)).
			Int("slow_transactions", len(txTransactions)).
			Int("fast_transactions", len(fastTxTransactions)).
			Int("slow_spent_txids", len(spentTxids)).
			Int("fast_spent_txids", len(fastSpentTxids)).
			Int("slow_spent_vouts", len(spentVouts)).
			Int("fast_spent_vouts", len(fastSpentVouts)).
			Msg("Transaction scan results")

		losses = append(losses, txLosses...)
		transactions = append(transactions, txTransactions...)
	}

	err = t.db.RecordBlockIndexResults(types.FromRPCBlock(block), types.FromRPCTxOutSetInfo(coinStats), blockStats, losses, transactions, spentTxids, spentVouts)
	if err != nil {
		return fmt.Errorf("failed to record block index results: %v", err)
	}

	elapsed := time.Since(start)
	log.Info().
		Int64("block_height", height).
		Stringer("elapsed", elapsed).
		Int("losses", len(losses)).
		Msg("Block processed")

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

		atLeastOneBurn := false
		for _, vout := range tx.Vout {
			if vout.ScriptPubKey.Type == "nulldata" {
				losses = append(losses, types.Loss{
					TxID:        tx.Txid,
					BlockHash:   blockhash,
					BlockHeight: blockHeight,
					Vout:        vout.N,
					Amount:      types.FromBTCString(types.BTCString(vout.Value)),
					BurnScript:  vout.ScriptPubKey.Hex,
				})

				atLeastOneBurn = true
			} else if t.bf.TestString(vout.ScriptPubKey.Hex) {
				exists, err := t.db.BurnScriptExists(vout.ScriptPubKey.Hex)
				if err != nil {
					log.Error().Err(err).Str("script", vout.ScriptPubKey.Hex).Msg("Failed to check if script exists")
					continue
				}

				log.Info().Str("script", vout.ScriptPubKey.Hex).Bool("exists", exists).Msg("Burn script identified")

				if exists {
					losses = append(losses, types.Loss{
						TxID:        tx.Txid,
						BlockHash:   blockhash,
						BlockHeight: blockHeight,
						Vout:        vout.N,
						Amount:      types.FromBTCString(types.BTCString(vout.Value)),
						BurnScript:  vout.ScriptPubKey.Hex,
					})

					atLeastOneBurn = true
				}
			}
		}

		if atLeastOneBurn {
			jsonTx, err := json.Marshal(tx)
			if err != nil {
				log.Error().Err(err).Str("txid", tx.Txid).Msg("Failed to marshal transaction")
				continue
			}

			txs = append(txs, types.Transaction{
				TxID:               tx.Txid,
				TransactionDetails: string(jsonTx),
				BlockHeight:        blockHeight,
				BlockHash:          blockhash,
			})
		}

	}

	return losses, txs, spentTxids, spentVouts
}

func (t *Tracker) fastScanTransactions(blockhash string, blockHeight int64, transactions []btypes.TransactionDetail) ([]types.Loss, []types.Transaction, []string, []int) {
	numCPUs := runtime.NumCPU()
	transactionCh := make(chan btypes.TransactionDetail, len(transactions))
	lossCh := make(chan types.Loss)
	txCh := make(chan types.Transaction)
	spentTxidVoutCh := make(chan [2]interface{})

	var wg sync.WaitGroup
	wg.Add(numCPUs)
	log.Info().Int("num_cpus", numCPUs).Msg("Starting fast transaction scan")

	for i := 0; i < numCPUs; i++ {
		go func() {
			defer wg.Done()
			for tx := range transactionCh {
				t.processTransaction(blockhash, blockHeight, tx, lossCh, txCh, spentTxidVoutCh)
			}
		}()
	}

	log.Info().
		Msg("Started runners")

	go func() {
		wg.Wait()
		log.Info().Msg("All runners done")
		close(lossCh)
		close(txCh)
		close(spentTxidVoutCh)
	}()

	log.Info().
		Msg("Sending transactions to channel")
	for _, tx := range transactions {
		transactionCh <- tx
	}
	close(transactionCh)
	log.Info().
		Msg("Transactions sent to channel")

	var losses []types.Loss
	var txs []types.Transaction
	var spentTxids []string
	var spentVouts []int

	for lossCh != nil || txCh != nil || spentTxidVoutCh != nil {
		select {
		case loss, ok := <-lossCh:
			if !ok {
				lossCh = nil
				continue
			}
			losses = append(losses, loss)
		case tx, ok := <-txCh:
			if !ok {
				txCh = nil
				continue
			}
			txs = append(txs, tx)
		case spent, ok := <-spentTxidVoutCh:
			if !ok {
				spentTxidVoutCh = nil
				continue
			}
			spentTxids = append(spentTxids, spent[0].(string))
			spentVouts = append(spentVouts, spent[1].(int))
		}
	}

	return losses, txs, spentTxids, spentVouts
}

func (t *Tracker) processTransaction(blockhash string, blockHeight int64, tx btypes.TransactionDetail, lossCh chan types.Loss, txCh chan types.Transaction, spentTxidVoutCh chan [2]interface{}) {
	var atLeastOneBurn bool

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
				spentTxidVoutCh <- [2]interface{}{vin.Txid, vin.Vout}
			}
		}
	}

	for _, vout := range tx.Vout {
		if vout.ScriptPubKey.Type == "nulldata" {
			lossCh <- types.Loss{
				TxID:        tx.Txid,
				BlockHash:   blockhash,
				BlockHeight: blockHeight,
				Vout:        vout.N,
				Amount:      types.FromBTCString(types.BTCString(vout.Value)),
				BurnScript:  vout.ScriptPubKey.Hex,
			}

			atLeastOneBurn = true
		} else if t.bf.TestString(vout.ScriptPubKey.Hex) {
			exists, err := t.db.BurnScriptExists(vout.ScriptPubKey.Hex)
			if err != nil {
				log.Error().Err(err).Str("script", vout.ScriptPubKey.Hex).Msg("Failed to check if script exists")
				continue
			}

			log.Info().Str("script", vout.ScriptPubKey.Hex).Bool("exists", exists).Msg("Burn script identified")

			if exists {
				lossCh <- types.Loss{
					TxID:        tx.Txid,
					BlockHash:   blockhash,
					BlockHeight: blockHeight,
					Vout:        vout.N,
					Amount:      types.FromBTCString(types.BTCString(vout.Value)),
					BurnScript:  vout.ScriptPubKey.Hex,
				}

				atLeastOneBurn = true
			}
		}
	}

	if atLeastOneBurn {
		jsonTx, err := json.Marshal(tx)
		if err != nil {
			log.Error().Err(err).Str("txid", tx.Txid).Msg("Failed to marshal transaction")
			return
		}

		txCh <- types.Transaction{
			TxID:               tx.Txid,
			TransactionDetails: string(jsonTx),
			BlockHeight:        blockHeight,
			BlockHash:          blockhash,
		}
	}
}
