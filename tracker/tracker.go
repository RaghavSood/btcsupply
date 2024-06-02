package tracker

import (
	"database/sql"
	"os"
	"time"

	"github.com/RaghavSood/btcsupply/bitcoinrpc"
	"github.com/RaghavSood/btcsupply/storage"
	"github.com/RaghavSood/btcsupply/types"
	zlog "github.com/rs/zerolog/log"
)

var log = zlog.With().Str("module", "tracker").Logger()

type Tracker struct {
	db     storage.Storage
	client *bitcoinrpc.RpcClient
}

func NewTracker(db storage.Storage) *Tracker {
	return &Tracker{
		db:     db,
		client: bitcoinrpc.NewRpcClient(os.Getenv("BITCOIND_HOST"), os.Getenv("BITCOIND_USER"), os.Getenv("BITCOIND_PASS")),
	}
}

func (t *Tracker) Run() {
	log.Info().Msg("Starting tracker")

	ticker := time.NewTicker(1 * time.Minute)
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
				Int("latest_block", latestBlock.BlockHeight).
				Int("current_block", info.Blocks).
				Msg("Checking for new blocks")

			for i := latestBlock.BlockHeight + 1; i <= info.Blocks; i++ {
				err = t.processBlock(i)
				if err != nil {
					log.Error().Err(err).Int("block_height", i).Msg("Failed to process block")
					break
				}
			}
		}
	}
}

func (t *Tracker) processBlock(height int) error {
	log.Info().Int("block_height", height).Msg("Processing block")
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
		Int("txouts", coinStats.Txouts).
		Float64("coinbase", coinStats.BlockInfo.Coinbase).
		Float64("prevout_spent", coinStats.BlockInfo.PrevoutSpent).
		Float64("new_outputs_ex_coinbase", coinStats.BlockInfo.NewOutputsExCoinbase).
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

	return nil
}
