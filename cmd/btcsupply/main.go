package main

import (
	"flag"
	"os"

	"github.com/RaghavSood/btcsupply/storage/sqlite"
	"github.com/RaghavSood/btcsupply/tracker"
	"github.com/RaghavSood/btcsupply/webui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var noindex bool

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: log.Output(zerolog.ConsoleWriter{Out: os.Stderr})})

	flag.BoolVar(&noindex, "noindex", false, "Don't index the blockchain, run in read-only mode")
	flag.Parse()
}

func main() {
	db, err := sqlite.NewSqliteBackend(noindex)
	defer db.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database")
	}

	if !noindex {
		tracker := tracker.NewTracker(db)
		go tracker.Run()
	} else {
		log.Info().Msg("Running in read-only mode, not indexing the blockchain")
	}

	webuiServer := webui.NewWebUI(db, noindex)
	webuiServer.Serve()
}
