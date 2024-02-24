package infrastructures

import (
	"sync"

	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	once     sync.Once        //nolint:gochecknoglobals //Singleton
	instance *TracerSingleton //nolint:gochecknoglobals //Singleton
)

type TracerSingleton struct {
	tracer trace.Tracer
	config SentryConfig
}

// SentryConfig holds the configuration for the sentry client.
type SentryConfig struct {
	Environment        string
	Release            string
	DSN                string
	Name               string
	TracesSampleRate   float64
	ProfilesSampleRate float64
	Debug              bool
	WithStacktrace     bool
}

func GetTracerInstance() *TracerSingleton {
	once.Do(func() {
		instance = &TracerSingleton{}
	})
	return instance
}

func NewTracerInstance(cfg SentryConfig) *TracerSingleton {
	return GetTracerInstance().WithConfig(cfg)
}

func (t *TracerSingleton) WithConfig(cfg SentryConfig) *TracerSingleton {
	t.config = cfg
	return t
}

func (t *TracerSingleton) Tracer() trace.Tracer {
	return t.tracer
}

func (t *TracerSingleton) ConnectTracer() error {
	t.tracer = otel.Tracer(t.config.Name)

	_ = sentry.Init(sentry.ClientOptions{
		Dsn:                t.config.DSN,
		Debug:              t.config.Debug,
		AttachStacktrace:   t.config.WithStacktrace,
		Environment:        t.config.Environment,
		Release:            t.config.Release,
		EnableTracing:      true,
		TracesSampleRate:   t.config.TracesSampleRate,
		ProfilesSampleRate: t.config.ProfilesSampleRate,
	})

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())

	return nil
}
