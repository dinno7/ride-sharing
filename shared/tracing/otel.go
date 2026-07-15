package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
	"go.opentelemetry.io/otel/trace"
)

type OTelConfig struct {
	ServiceName string
	Environment string
	ExporterURL string
}

func InitTracer(cfg OTelConfig) (*sdktrace.TracerProvider, error) {
	ctx := context.Background()
	// TODO: Exporter
	traceExporter, err := newExporter(ctx, cfg.ExporterURL)
	if err != nil {
		return nil, fmt.Errorf("failed to init trace exporter: %w", err)
	}

	// TODO: Trace provider
	traceProvider, err := newTraceProvider(ctx, &cfg, traceExporter)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(traceProvider)

	// TODO: Propagator
	propagator := newPropagator()
	otel.SetTextMapPropagator(propagator)

	return traceProvider, nil
}

func GetTracer(name string) trace.Tracer {
	return otel.GetTracerProvider().Tracer(name)
}

func newTraceProvider(
	ctx context.Context,
	cfg *OTelConfig,
	exporter sdktrace.SpanExporter,
) (*sdktrace.TracerProvider, error) {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.DeploymentEnvironmentNameKey.String(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace provider's resource: %w", err)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	return traceProvider, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newExporter(ctx context.Context, endpoint string) (sdktrace.SpanExporter, error) {
	return otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint))
}
