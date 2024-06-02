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
		}
	}
}
