package logger

import (
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/rs/zerolog"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"os"
)

func CreateNewRelicLogger() {
	wr := diode.NewWriter(newRollingFile(), config.DiodeLoggerSize, config.DiodeLoggerTime, func(missed int) {
		log.Printf("Logger Dropped %d messages", missed)
	})
	writer := zerologWriter.New(wr, infrastructures.GetRelicApp())
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		writer = zerologWriter.New(prettyLogger(), infrastructures.GetRelicApp())

	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "none":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	default:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}

	writer.DebugLogging(true)

	log.Logger = zerolog.New(writer)

}
