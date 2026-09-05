package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func InitLogger(env string) {
	if env == "development" {
		Log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}).
			With().Timestamp().Logger()
	} else {
		Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}
