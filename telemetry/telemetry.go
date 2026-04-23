package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Enabled     bool    `env:"OTEL_ENABLED" envDefault:"false"`
	Endpoint    string  `env:"OTEL_EXPORTER_OTLP_ENDPOINT" envDefault:"localhost:4317"`
	ServiceName string  `env:"OTEL_SERVICE_NAME" envDefault:"appfr"`
	SampleRatio float64 `env:"OTEL_TRACES_SAMPLE_RATIO" envDefault:"1.0"`
	Insecure    bool    `env:"OTEL_EXPORTER_OTLP_INSECURE" envDefault:"true"`
}

type Telemetry struct {
	tracerProvider *sdktrace.TracerProvider
	config         Config
}

func New(ctx context.Context, cfg Config) (*Telemetry, error) {
	if !cfg.Enabled {
		return &Telemetry{config: cfg}, nil
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	}
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SampleRatio))

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Telemetry{
		tracerProvider: tp,
		config:         cfg,
	}, nil
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	if t.tracerProvider == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return t.tracerProvider.ForceFlush(ctx)
}

func (t *Telemetry) IsEnabled() bool {
	return t.config.Enabled
}

func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
