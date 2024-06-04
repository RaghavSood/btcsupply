package electrum

import (
	"context"
	"fmt"
	"time"

	"github.com/RaghavSood/btcsupply/btclogger"
	electrumx "github.com/checksum0/go-electrum/electrum"
)

var log = btclogger.NewLogger("electrum")

type Electrum struct {
	client *electrumx.Client
	Height int64
}

func NewElectrum() (*Electrum, error) {
	client, err := electrumx.NewClientTCP(context.Background(), "electrum3.bluewallet.io:50001")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to connect to Electrum server")

		return nil, err
	}

	e := &Electrum{
		client: client,
	}

	// Ping the server every 60 seconds to keep the connection alive
	go e.Ping()

	// Making sure we declare to the server what protocol we want to use
	serverVer, protocolVer, err := e.client.ServerVersion(context.TODO())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get server version")
	}

	// Subscribe to header updates so we can keep track of the current height
	go e.MonitorHeaders()

	log.Info().
		Str("server_version", serverVer).
		Str("protocol_version", protocolVer).
		Msg("Connected to Electrum server")

	return e, nil
}

func (e *Electrum) Ping() {
	for {
		if err := e.client.Ping(context.TODO()); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to ping server")
		}
		time.Sleep(10 * time.Second)
	}
}

func (e *Electrum) MonitorHeaders() {
	sub, err := e.client.SubscribeHeaders(context.Background())
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to subscribe to headers")
	}

	log.Info().Msg("Subscribed to headers")

	for {
		select {
		case header := <-sub:
			height := int64(header.Height)
			log.Info().
				Int64("height", height).
				Msg("New block height")
			e.Height = height
		}
	}
}

func (e *Electrum) GetScriptHistory(script string) ([]string, error) {
	electrumHash, err := ScriptToElectrumScript(script)
	if err != nil {
		return nil, fmt.Errorf("failed to convert script to electrum hash: %w", err)
	}

	history, err := e.client.GetHistory(context.TODO(), electrumHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get script history: %w", err)
	}

	var txids []string
	for _, entry := range history {
		log.Debug().
			Str("txid", entry.Hash).
			Int32("height", entry.Height).
			Uint32("fee", entry.Fee).
			Msg("Found transaction")
		txids = append(txids, entry.Hash)
	}

	return txids, nil
}

func (e *Electrum) GetScriptUnspents(script string) ([]string, []int64, error) {
	electrumHash, err := ScriptToElectrumScript(script)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert script to electrum hash: %w", err)
	}

	unspents, err := e.client.ListUnspent(context.TODO(), electrumHash)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get script unspents: %w", err)
	}

	var unspentTxids []string
	var unspentHeights []int64
	seen := make(map[string]bool)
	for _, entry := range unspents {
		log.Debug().
			Str("txid", entry.Hash).
			Uint32("height", entry.Height).
			Msg("Found unspent transaction")

		if seen[entry.Hash] {
			continue
		}

		// Block heights 0 and -1 indicate unconfirmed transactions
		// These will be picked up by our node indexing normally once they
		// are confirmed
		if entry.Height > 0 {
			unspentTxids = append(unspentTxids, entry.Hash)
			unspentHeights = append(unspentHeights, int64(entry.Height))
		}

		seen[entry.Hash] = true
	}

	return unspentTxids, unspentHeights, nil
}
