// Package trace provides a factory method to create a tracer provider.
//
// As your infrastructure grows, it becomes important to be able to trace a
// request, as it travels through multiple services and back to the user.
package trace

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Default values
const defaultProbability = 0.1 // 10%

// Supported exporters.
const (
	ZipKinExporter = "zipkin"
	GRPCExporter   = "otlpgrpc"
)

// Config is the gtrace configuration.
type Config struct {
	Exporter           string
	ServiceName        string
	Probability        float64
	ReporterURI        string
	APIKey             string
	EnableTLS          bool
	ReconnectionPeriod time.Duration
	ConnectionTimeout  time.Duration
}

// StartTracing initialises the tracer provider with the given trace configuration.
func StartTracing(cfg Config) (*sdktrace.TracerProvider, error) {
	if cfg.Exporter == "" {
		return nil, fmt.Errorf("exporter cannot be empty")
	}

	if cfg.ReporterURI == "" {
		return nil, fmt.Errorf("endpointURL cannot be empty")
	}

	if cfg.Probability == 0 {
		cfg.Probability = defaultProbability
	}

	exporter, err := newExporter(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating span exporter: %w", err)
	}

	// Create trace provider.
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.Probability)),
		sdktrace.WithBatcher(exporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
		),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(cfg.ServiceName),
				attribute.String("exporter", cfg.Exporter),
			),
		),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return traceProvider, nil
}

func newExporter(cfg Config) (sdktrace.SpanExporter, error) {
	switch cfg.Exporter {
	case ZipKinExporter:
		return zipkin.New(cfg.ReporterURI)
	case GRPCExporter:
		return newGRPCExporter(cfg)
	default:
		return nil, fmt.Errorf("exporter [%s] not supported", cfg.Exporter)
	}
}

func newGRPCExporter(cfg Config) (*otlptrace.Exporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.ReporterURI),
		otlptracegrpc.WithReconnectionPeriod(cfg.ReconnectionPeriod),
		otlptracegrpc.WithTimeout(cfg.ConnectionTimeout),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	}

	if cfg.EnableTLS {
		cred := credentials.NewClientTLSFromCert(nil, "")
		opts = append(opts, otlptracegrpc.WithTLSCredentials(cred))
	}

	if cfg.APIKey != "" {
		headers := map[string]string{"api-key": cfg.APIKey}
		opts = append(opts, otlptracegrpc.WithHeaders(headers))
	}

	exp, err := otlptrace.New(context.Background(), otlptracegrpc.NewClient(opts...))
	if err != nil {
		return nil, fmt.Errorf("craeting otlptracegrpc client: %w", err)
	}

	return exp, nil
}
