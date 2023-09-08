package infrastructures

import (
	"context"
	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/memnix/memnix-rest/config"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	fiberTracer = otel.Tracer("fiber-server")
)

func InitTracer(cfg config.SentryConfigStruct) error {

	initSentry(cfg)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())

	return nil
}

func ShutdownTracer() error {
	return nil
}

func GetFiberTracer() trace.Tracer {
	return fiberTracer
}

func initSentry(cfg config.SentryConfigStruct) {
	otelzap.Ctx(context.Background()).Info("Initializing Sentry :", zap.Float64("traces_sample_rate", cfg.TracesSampleRate), zap.Float64("profiles_sample_rate", cfg.ProfilesSampleRate))
	_ = sentry.Init(sentry.ClientOptions{
		Dsn:                cfg.DSN,
		Debug:              cfg.Debug,
		AttachStacktrace:   true,
		Environment:        cfg.Environment,
		Release:            cfg.Release,
		EnableTracing:      true,
		TracesSampleRate:   cfg.TracesSampleRate,
		ProfilesSampleRate: cfg.ProfilesSampleRate,
	})
}
