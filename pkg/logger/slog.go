package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	logger *Logger   //nolint:gochecknoglobals //Singleton
	once   sync.Once //nolint:gochecknoglobals //Singleton
)

type Logger struct {
	logLevel slog.Level
}

func GetLogger() *Logger {
	once.Do(func() {
		logger = &Logger{
			logLevel: slog.LevelInfo,
		}
	})
	return logger
}

func (l *Logger) SetLogLevel(level slog.Level) *Logger {
	l.logLevel = level
	return l
}

func (l *Logger) GetLogLevel() slog.Level {
	return l.logLevel
}

func (l *Logger) CreateGlobalHandler() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     l.logLevel,
		AddSource: false,
	})

	logger := slog.New(handler)

	slog.SetLogLoggerLevel(l.logLevel)

	slog.SetDefault(logger)
}
