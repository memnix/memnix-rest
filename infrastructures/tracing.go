package infrastructures

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	fiberTracer = otel.Tracer("fiber-server")
	tp          *sdktrace.TracerProvider
)

func buildEndpoint(url string) string {
	var sb strings.Builder
	sb.WriteString(url)
	sb.WriteString("/api/traces")
	return sb.String()
}

func InitTracer(url string) error {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(buildEndpoint(url))))
	if err != nil {
		return err
	}

	tp = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("memnix-backend"),
			)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return nil
}

func ShutdownTracer() error {
	return tp.Shutdown(context.Background())
}

func GetTracer() *sdktrace.TracerProvider {
	return tp
}

func GetFiberTracer() trace.Tracer {
	return fiberTracer
}
