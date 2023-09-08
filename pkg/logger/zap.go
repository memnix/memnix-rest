package logger

import (
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/memnix/memnix-rest/config"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ConvertSentryLevel converts a zapcore.Level to a sentry.Level
func ConvertSentryLevel(level zapcore.Level) sentry.Level {
	switch level {
	case zapcore.InvalidLevel:
		return sentry.LevelFatal
	case zapcore.DebugLevel:
		return sentry.LevelDebug
	case zapcore.InfoLevel:
		return sentry.LevelInfo
	case zapcore.WarnLevel:
		return sentry.LevelWarning
	case zapcore.ErrorLevel:
		return sentry.LevelError
	case zapcore.DPanicLevel:
		return sentry.LevelFatal
	case zapcore.PanicLevel:
		return sentry.LevelFatal
	case zapcore.FatalLevel:
		return sentry.LevelFatal
	default:
		return sentry.LevelFatal
	}
}

func CreateZapLogger() (*otelzap.Logger, func()) {
	zapLogger, _ := zap.NewProduction(zap.Hooks(func(entry zapcore.Entry) error {
		defer sentry.Flush(config.SentryFlushTimeout)
		if entry.Level <= zapcore.InfoLevel {
			return nil
		}
		msg := strings.Builder{}
		msg.WriteString(entry.Message)
		msg.WriteString(" ")
		msg.WriteString(entry.Caller.TrimmedPath())
		event := sentry.NewEvent()
		event.Message = msg.String()
		event.Level = ConvertSentryLevel(entry.Level)
		sentry.CaptureEvent(event)
		return nil
	}))

	logger := otelzap.New(zapLogger, otelzap.WithTraceIDField(false))
	undo := otelzap.ReplaceGlobals(logger)

	return logger, undo
}
