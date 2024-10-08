package trc

import (
	"context"
	"packagelock/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

var Tracer *trace.TracerProvider

// Initializes the TracerProvider
func InitTracer() (*trace.TracerProvider, error) {
	ctx := context.Background()

	// Create an OTLP gRPC exporter
	conn, err := grpc.DialContext(ctx, "localhost:4317", grpc.WithInsecure()) // change to your exporter endpoint
	if err != nil {
		logger.Logger.Errorf("failed to create the OTLP trace exporter: %w", err)
		return nil, err
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		logger.Logger.Errorf("failed to create the OTLP trace exporter: %w", err)
		return nil, err
	}

	// Create a new TracerProvider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("PackageLock"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
