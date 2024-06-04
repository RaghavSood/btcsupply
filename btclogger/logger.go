package btclogger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(module string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	return zerolog.
		New(consoleWriter).
		With().
		Timestamp().
		Str("module", module).
		Logger()
}
