package tracing

import (
	"context"
	"fmt"
	"os"

	// For setting span status

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewTracerProvider creates and configures a new TracerProvider
// It looks for two env flags:
//   - 'TRACING_ENABLED' -> enables tracing, if not set/false uses noop-tracer
//   - 'TRACING_JAEGER_URL' -> sets custom jaeger url, if not set uses default localhost
func NewTracerProvider(logger *zap.Logger) (*sdktrace.TracerProvider, error) {
	tracingEnabled := os.Getenv("TRACING_ENABLED") == "true"
	if !tracingEnabled {
		logger.Info("Tracing is disabled via ENV-FLAG ('TRACING_ENABLED'). Using NoopTracerProvider.")
		return sdktrace.NewTracerProvider(), nil // Returns a TracerProvider with no exporters
	}
	logger.Info("Tracing is enabled via ENV-FLAG ('TRACING_ENABLED'). Using TracerProvider.")

	var exporter *jaeger.Exporter
	var exporterErr error

	jaegerAddress := os.Getenv("TRACING_JAEGER_URL")

	if jaegerAddress != "" {

		// Configure the Jaeger exporter with custom URL from ENV-FLAG
		exporter, exporterErr = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerAddress)))
		if exporterErr != nil {
			logger.Error("Failed to create Jaeger exporter", zap.Error(exporterErr))
			return nil, fmt.Errorf("failed to create Jaeger exporter: %w", exporterErr)
		}
		logger.Info("Configured Jaeger exporter with custom url from ENV-FLAG.")

	} else {

		// Configure the Jaeger exporter with default url
		exporter, exporterErr = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
		if exporterErr != nil {
			logger.Error("Failed to create Jaeger exporter", zap.Error(exporterErr))
			return nil, fmt.Errorf("failed to create Jaeger exporter: %w", exporterErr)
		}
	}

	// Create the TracerProvider with batching and resource attributes
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("PackageLock"),
		)),
	)

	// Set the global TracerProvider
	otel.SetTracerProvider(tp)

	logger.Info("Tracing is enabled and TracerProvider is configured.")
	return tp, nil
}

// NewTracer provides a trace.Tracer instance
func NewTracer(tp *sdktrace.TracerProvider, logger *zap.Logger) trace.Tracer {
	tracingEnabled := os.Getenv("TRACING_ENABLED") == "true"
	if !tracingEnabled {
		logger.Info("Tracing is disabled via ENV-FLAG ('TRACING_ENABLED'). Using Nooptrace.Tracer.")
		return otel.GetTracerProvider().Tracer("noop-tracer")
	}
	logger.Info("Tracing is enabled via ENV-FLAG ('TRACING_ENABLED'). Using trace.Tracer.")
	return tp.Tracer("PackageLock")
}

// Module is the FX module for tracing
var Module = fx.Options(
	fx.Provide(
		NewTracerProvider, // Provides *sdktrace.TracerProvider
		NewTracer,         // Provides trace.Tracer
	),
	fx.Invoke(func(lc fx.Lifecycle, tp *sdktrace.TracerProvider, logger *zap.Logger) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				if err := tp.Shutdown(ctx); err != nil {
					logger.Error("Error shutting down TracerProvider", zap.Error(err))
					return err
				}
				logger.Info("TracerProvider shutdown successfully.")
				return nil
			},
		})
	}),
)
