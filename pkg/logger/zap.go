package logger

import (
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func CreateZapLogger() (*otelzap.Logger, func()) {
	zapLogger, _ := zap.NewProduction()

	logger := otelzap.New(zapLogger)
	undo := otelzap.ReplaceGlobals(logger)

	return logger, undo
}
