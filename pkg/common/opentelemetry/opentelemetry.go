package opentelemetry

import (
	"context"
	"errors"
	"net/http"
	"time"

	cfgcommon "github.com/durianpay/dpay-common/config"
	"github.com/durianpay/dpay-common/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/sdk/metric"
	sdkResource "go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

const (
	// ScopeName is the instrumentation scope name.
	Name         = "github.com/durianpay/dpay-common/opentelemetry" // change as we need later
	Version      = "0.0.1"
	RequestIDKey = "trace.request.id"
)

var (
	appName    string
	appVersion string
	appEnv     cfgcommon.Environment
)

func init() {
	appName = cfgcommon.AppName()
	appEnv = cfgcommon.Env()
	appVersion := "" // TODO can add via env

	if appName == "" {
		appName = "disbursement_service"
	}

	if appVersion == "" {
		appVersion = "1.0.0"
	}
}

type closeFn func()

func InitOTelTrace(
	ctx context.Context,
	enableOpenTelemtry bool,
) (closeFn, error) {
	var err error
	// Configure and create an OTLP exporter to send traces to an OLTP collector.

	resource, err := newResource(ResourceOptions{
		ServiceName:       appName,
		ServiceVersion:    appVersion,
		ServiceInstanceID: "", // TODO: get service instance id information
		Deployment: ResourceDeploymentOptions{
			Environment: string(appEnv),
		},
	})
	if err != nil {
		logger.Errorf(ctx, "Error creating tracer resource")
		return disableOtel(), err
	}

	tracerProvider, err := setupTracerProvider(TracerProviderOptions{
		WithInsecure: true,
	}, resource)
	if err != nil {
		logger.Errorf(ctx, "Error creating tracer resource")
		return disableOtel(), err
	}

	meterProvider, err := setupMeterProvider(resource)
	if err != nil {
		logger.Errorf(ctx, "Error creating tracer resource")
		return disableOtel(), err
	}

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(propagator)

	return func() {
		go func() {
			<-ctx.Done()

			defer otel.SetTracerProvider(noop.NewTracerProvider())

			err = tracerProvider.Shutdown(context.TODO())
			if nil != err && !errors.Is(err, context.Canceled) {
				logger.Errorf(context.TODO(), "Error shutting down tracer provider")
			}

			err = meterProvider.Shutdown(context.TODO())
			if nil != err && !errors.Is(err, context.Canceled) {
				logger.Errorf(context.TODO(), "Error shutting down meter provider")
			}
		}()
	}, nil
}

func GetTracer() trace.Tracer {
	tracer := otel.GetTracerProvider().Tracer(Name, trace.WithInstrumentationVersion(Version))
	return tracer
}

// GetSpan retrieves the current trace span from the given context.
// If no span is found, it returns NoOp span.
func GetSpan(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// GetOrCreateSpan retrieves the current trace span from the given context.
// If no span is found, it creates a new span with the specified name and options.
// It returns the updated context and the obtained or newly created span.
func GetOrCreateSpan(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	span := GetSpan(ctx)
	if !span.IsRecording() {
		ctx, span = StartSpan(ctx, spanName, opts...)
	}

	return ctx, span
}

// StartSpan creates a new trace span with the specified name and options using the current tracer.
// It returns the updated context and the newly created span.
func StartSpan(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	tracer := GetTracer()

	opts = append([]trace.SpanStartOption{
		trace.WithTimestamp(time.Now().UTC()),
	}, opts...)

	return tracer.Start(ctx, spanName, opts...)
}

func WrapHTTPHandler(h http.Handler, operation string, opts ...otelhttp.Option) http.Handler {
	return otelhttp.NewHandler(h, operation, opts...)
}

func MetricHandler() http.Handler {
	return promhttp.Handler()
}

func disableOtel() closeFn {
	otel.SetTracerProvider(noop.NewTracerProvider())

	return closeFn(func() {})
}

func setupTracerProvider(
	opt TracerProviderOptions,
	resource *sdkResource.Resource,
) (*sdktrace.TracerProvider, error) {
	tracerProvider, err := newTracerProvider(opt, resource)
	if err != nil {
		return tracerProvider, err
	}

	otel.SetTracerProvider(tracerProvider)

	return tracerProvider, nil
}

func setupMeterProvider(r *sdkResource.Resource) (*metric.MeterProvider, error) {
	meterProvider, err := newMeterProvider(r)
	if err != nil {
		return meterProvider, err
	}

	otel.SetMeterProvider(meterProvider)

	return meterProvider, nil
}
