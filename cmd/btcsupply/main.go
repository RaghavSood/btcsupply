package main

import (
	"os"

	"github.com/RaghavSood/btcsupply/storage/sqlite"
	"github.com/RaghavSood/btcsupply/webui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: log.Output(zerolog.ConsoleWriter{Out: os.Stderr})})
}

func main() {
	db, err := sqlite.NewSqliteBackend()
	db.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database")
	}

	webuiServer := webui.NewWebUI(db)
	webuiServer.Serve()
}
