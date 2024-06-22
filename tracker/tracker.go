package tracker

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
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

func (t *Tracker) getBlockStats(height int64, blockStatsCh chan btypes.BlockStats, errCh chan error, wg *sync.WaitGroup) {
	startBlockStats := time.Now()
	defer wg.Done()
	blockStats, err := t.client.GetBlockStats(height)
	if err != nil {
		errCh <- err
		return
	}

	blockStatsCh <- blockStats
	elapsedBlockStats := time.Since(startBlockStats)
	log.Info().Stringer("block_stats_elapsed", elapsedBlockStats).Int64("height", height).Msg("Block stats fetched")
}

func (t *Tracker) getCoinStats(height int64, txOutSetInfoCh chan btypes.TxOutSetInfo, errCh chan error, wg *sync.WaitGroup) {
	startCoinStats := time.Now()
	defer wg.Done()
	coinStats, err := t.client.GetTxOutSetInfo("muhash", height, true)
	if err != nil {
		errCh <- err
		return
	}

	txOutSetInfoCh <- coinStats
	elapsedCoinStats := time.Since(startCoinStats)
	log.Info().Stringer("coin_stats_elapsed", elapsedCoinStats).Int64("height", height).Msg("Coin stats fetched")
}

func (t *Tracker) getBlock(height int64, blockCh chan btypes.Block, errCh chan error, wg *sync.WaitGroup) {
	startBlockTime := time.Now()
	defer wg.Done()
	hash, err := t.client.GetBlockHash(height)
	if err != nil {
		errCh <- err
		return
	}

	block, err := t.client.GetBlock(hash)
	if err != nil {
		errCh <- err
		return
	}

	blockCh <- block
	elapsedBlockTime := time.Since(startBlockTime)
	log.Info().Stringer("block_time_elapsed", elapsedBlockTime).Int64("height", height).Msg("Block fetched")
}

func (t *Tracker) processBlock(height int64) error {
	log.Info().Int64("block_height", height).Msg("Processing block")
	start := time.Now()

	blockStatsCh := make(chan btypes.BlockStats)
	txOutSetInfoCh := make(chan btypes.TxOutSetInfo)
	blockCh := make(chan btypes.Block)
	errCh := make(chan error, 3)

	var wg sync.WaitGroup
	wg.Add(3)

	go t.getBlockStats(height, blockStatsCh, errCh, &wg)
	go t.getCoinStats(height, txOutSetInfoCh, errCh, &wg)
	go t.getBlock(height, blockCh, errCh, &wg)

	go func() {
		wg.Wait()
		close(blockStatsCh)
		close(txOutSetInfoCh)
		close(blockCh)
		close(errCh)
	}()

	var blockStats btypes.BlockStats
	var coinStats btypes.TxOutSetInfo
	var block btypes.Block

	for blockStatsCh != nil || txOutSetInfoCh != nil || blockCh != nil {
		select {
		case stats, ok := <-blockStatsCh:
			if !ok {
				blockStatsCh = nil
				continue
			}
			blockStats = stats
		case coin, ok := <-txOutSetInfoCh:
			if !ok {
				txOutSetInfoCh = nil
				continue
			}
			coinStats = coin
		case b, ok := <-blockCh:
			if !ok {
				blockCh = nil
				continue
			}
			block = b
		case err := <-errCh:
			if err != nil {
				return fmt.Errorf("failed to fetch block data: %v", err)
			}
		}
	}

	log.Info().
		Int64("subsidy", blockStats.Subsidy).
		Str("hash", blockStats.Blockhash).
		Int64("totalfee", blockStats.Totalfee).
		Int64("height", blockStats.Height).
		Msg("Block stats")

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

		log.Info().
			Int64("block_height", height).
			Stringer("elapsed", elapsed).
			Int("losses", len(txLosses)).
			Int("transactions", len(txTransactions)).
			Int("spent_txids", len(spentTxids)).
			Int("spent_vouts", len(spentVouts)).
			Msg("Block transactions scanned")

		losses = append(losses, txLosses...)
		transactions = append(transactions, txTransactions...)
	}

	startRecord := time.Now()
	err := t.db.RecordBlockIndexResults(types.FromRPCBlock(block), types.FromRPCTxOutSetInfo(coinStats), blockStats, losses, transactions, spentTxids, spentVouts)
	if err != nil {
		return fmt.Errorf("failed to record block index results: %v", err)
	}
	endRecord := time.Since(startRecord)

	elapsed := time.Since(start)
	log.Info().
		Stringer("record_elapsed", endRecord).
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
