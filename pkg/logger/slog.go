package logger

import (
	"log/slog"
	"os"
	"runtime/debug"
)

type Logger struct {
	logLevel slog.Level
}

func NewLogger() *Logger {
	return &Logger{
		logLevel: slog.LevelInfo,
	}
}

func (l *Logger) SetLogLevel(level slog.Level) {
	l.logLevel = level
}

func (l *Logger) GetLogLevel() slog.Level {
	return l.logLevel
}

func (l *Logger) CreateSlogHandler() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     l.logLevel,
		AddSource: l.logLevel == slog.LevelDebug,
	})

	buildInfo, _ := debug.ReadBuildInfo()

	logger := slog.New(handler)

	child := logger.With(
		slog.Group("program_info",
			slog.Int("pid", os.Getpid()),
			slog.String("go_version", buildInfo.GoVersion),
		))

	slog.SetDefault(child)
}
