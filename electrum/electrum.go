package electrum

import (
	"context"
	"time"

	electrumx "github.com/checksum0/go-electrum/electrum"
	zlog "github.com/rs/zerolog/log"
)

var log = zlog.With().Str("module", "electrum").Logger()

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
	}, err
}
