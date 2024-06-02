package tracker

import (
	"os"
	"time"

	"github.com/RaghavSood/btcsupply/bitcoinrpc"
	"github.com/RaghavSood/btcsupply/storage"
	"github.com/rs/zerolog/log"
)

var logger = log.With().Str("module", "tracker").Logger()

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
	logger.Info().Msg("Starting tracker")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Debug().Msg("Checking for new blocks")
			info, err := t.client.GetBlockchainInfo()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to get blockchain info")
				continue
			}

			latestBlock, err := t.db.GetLatestBlock()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to get latest block")
				continue
			}

			logger.Info().
				Int("latest_block", latestBlock.BlockHeight).
				Int("current_block", info.Blocks).
				Msg("Checking for new blocks")
		}
	}
}
