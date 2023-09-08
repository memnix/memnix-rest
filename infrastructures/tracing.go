package infrastructures

import (
	"context"
	"github.com/memnix/memnix-rest/config"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
)

var (
	otlpExporter *otlptrace.Exporter
	fiberTracer  = otel.Tracer("fiber-server")
)

func InitTracer(cfg config.TracingConfigStruct) error {

	var secureOption otlptracegrpc.Option

	if cfg.InsecureMode {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(cfg.OtelEndpoint),
		),
	)

	if err != nil {
		return errors.Wrap(err, "failed to create exporter")
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create resource")
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	otlpExporter = exporter

	return nil
}

func ShutdownTracer() error {
	return otlpExporter.Shutdown(context.Background())
}

func GetTracer() *sdktrace.TracerProvider {
	return otel.GetTracerProvider().(*sdktrace.TracerProvider)
}

func GetFiberTracer() trace.Tracer {
	return fiberTracer
}
