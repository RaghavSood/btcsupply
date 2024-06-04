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
}

func NewElectrum() (*Electrum, error) {
	client, err := electrumx.NewClientTCP(context.Background(), "electrum3.bluewallet.io:50001")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to connect to Electrum server")

		return nil, err
	}

	// Ping the server every 60 seconds to keep the connection alive
	go func() {
		for {
			if err := client.Ping(context.TODO()); err != nil {
				log.Fatal().
					Err(err).
					Msg("Failed to ping server")
			}
			time.Sleep(60 * time.Second)
		}
	}()

	// Making sure we declare to the server what protocol we want to use
	serverVer, protocolVer, err := client.ServerVersion(context.TODO())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get server version")
	}

	log.Info().
		Str("server_version", serverVer).
		Str("protocol_version", protocolVer).
		Msg("Connected to Electrum server")

	return &Electrum{
		client: client,
	}, nil
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
