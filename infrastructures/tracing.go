package infrastructures

import (
	"context"
	"log/slog"
	"sync"

	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/gofiber/fiber/v2/log"
	"github.com/memnix/memnix-rest/config"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	once     sync.Once    //nolint:gochecknoglobals // Singleton
	instance trace.Tracer //nolint:gochecknoglobals // Singleton
)

func GetTracerInstance() trace.Tracer {
	once.Do(func() {
		instance = otel.Tracer("fiber-server")
	})
	return instance
}

func InitTracer(cfg config.SentryConfig) error {
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

func initSentry(cfg config.SentryConfig) {
	log.WithContext(context.Background()).Info(
		"Initializing Sentry :", slog.Float64("traces_sample_rate", cfg.TracesSampleRate), slog.Float64(
			"profiles_sample_rate", cfg.ProfilesSampleRate))
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
