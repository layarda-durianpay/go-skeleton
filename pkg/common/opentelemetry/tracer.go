package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

type TracerProviderOptions struct {
	WithInsecure      bool
	CollectorEndpoint string
}

func newTracerProvider(
	opt TracerProviderOptions,
	resourceObj *resource.Resource,
) (*trace.TracerProvider, error) {
	var opts []otlptracegrpc.Option

	if opt.CollectorEndpoint != "" {
		opts = append(opts, otlptracegrpc.WithEndpoint(opt.CollectorEndpoint))
	}

	if opt.WithInsecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resourceObj),
	)

	return traceProvider, nil
}
