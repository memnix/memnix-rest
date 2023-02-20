package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/memnix/memnix-rest/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func prettyLogger() zerolog.ConsoleWriter {
	return zerolog.ConsoleWriter{Out: os.Stderr,
		TimeFormat: zerolog.TimeFormatUnix,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
	}
}

func CreateLogger() {
	wr := diode.NewWriter(newRollingFile(), config.DiodeLoggerSize, config.DiodeLoggerTime, func(missed int) {
		log.Printf("Logger Dropped %d messages", missed)
	})
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(prettyLogger()).With().Caller().Logger()
		return
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
	log.Logger = log.Output(wr).With().Caller().Logger()
}

func newRollingFile() io.Writer {
	return &lumberjack.Logger{
		Filename:   path.Join("./logs", "logs.log"),
		MaxBackups: config.MaxBackupLogFiles, // files
		MaxSize:    config.MaxSizeLogFiles,   // megabytes
	}
}
