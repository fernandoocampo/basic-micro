package telemetry

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type ShutdownFunc func(context.Context)

// OtelSDKSetup contains data needed by telemetry implementation.
type OtelSDKSetup struct {
	ServiceName    string
	ServiceVersion string
	Logger         *slog.Logger
}

func NewOtelSDK(ctx context.Context, otelSetup OtelSDKSetup) (shutdown ShutdownFunc, err error) {
	var shutdownFuncs []func(context.Context) error

	logger := otelSetup.Logger.With(
		slog.Group("telemetry",
			slog.String("provider", "otel"),
		),
	)

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) {
		for _, fn := range shutdownFuncs {
			fn(ctx)
		}
		shutdownFuncs = nil
	}

	// Set up resource.
	logger.Info("setting up resource")
	res, err := newResource(otelSetup)
	if err != nil {
		logger.Error("creating new resource", "error", err)
		shutdown(ctx)
		return
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	logger.Info("setting up trace provider")
	tracerProvider, err := newTraceProvider(res)
	if err != nil {
		logger.Error("creating new trace provider", "error", err)
		shutdown(ctx)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	logger.Info("setting up meter provider")
	meterProvider, err := newMeterProvider(res)
	if err != nil {
		logger.Error("creating new meter provider", "error", err)
		shutdown(ctx)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func newResource(otelSetup OtelSDKSetup) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(otelSetup.ServiceName),
			semconv.ServiceVersion(otelSetup.ServiceVersion),
		))
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

func newMeterProvider(res *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),
	)
	return meterProvider, nil
}
