package observability

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"erp-job/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	tracer trace.Tracer
	meter  metric.Meter

	runsTotal           metric.Int64Counter
	runDuration         metric.Float64Histogram
	recordsFetchedTotal metric.Int64Counter
	recordsPostedTotal  metric.Int64Counter
	retriesTotal        metric.Int64Counter
	failuresTotal       metric.Int64Counter
	checkpointLag       metric.Int64Histogram
	lastSuccessGauge    metric.Int64ObservableGauge

	lastSuccessUnix atomic.Int64
}

type HTTPAttempt struct {
	Endpoint    string
	Attempt     int
	StatusCode  int
	Error       error
	WillRetry   bool
	Duration    time.Duration
	EndpointTag string
}

type attemptObserverKey struct{}
type runIDKey struct{}

type AttemptObserver func(HTTPAttempt)

func New(ctx context.Context, cfg config.OTel) (*Telemetry, func(context.Context) error, error) {
	telemetry := &Telemetry{}

	if !cfg.Enabled {
		if err := telemetry.initInstruments(); err != nil {
			return nil, nil, err
		}
		return telemetry, func(context.Context) error { return nil }, nil
	}

	traceOpts, metricOpts, err := exporterOptions(cfg)
	if err != nil {
		return nil, nil, err
	}

	traceExporter, err := otlptracehttp.New(ctx, traceOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("create trace exporter: %w", err)
	}

	metricExporter, err := otlpmetrichttp.New(ctx, metricOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("create metric exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			attribute.String("deployment.environment", cfg.Environment),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create telemetry resource: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(5*time.Second))),
		sdkmetric.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)

	if err := telemetry.initInstruments(); err != nil {
		return nil, nil, err
	}

	shutdown := func(ctx context.Context) error {
		var shutdownErr error
		if err := meterProvider.Shutdown(ctx); err != nil {
			shutdownErr = err
		}
		if err := tracerProvider.Shutdown(ctx); err != nil && shutdownErr == nil {
			shutdownErr = err
		}
		return shutdownErr
	}

	return telemetry, shutdown, nil
}

func (t *Telemetry) initInstruments() error {
	t.tracer = otel.Tracer("erp-job")
	t.meter = otel.Meter("erp-job")

	var err error
	if t.runsTotal, err = t.meter.Int64Counter("erp_job_runs_total"); err != nil {
		return fmt.Errorf("create runs_total counter: %w", err)
	}
	if t.runDuration, err = t.meter.Float64Histogram("erp_job_run_duration_seconds"); err != nil {
		return fmt.Errorf("create run_duration histogram: %w", err)
	}
	if t.recordsFetchedTotal, err = t.meter.Int64Counter("erp_job_records_fetched_total"); err != nil {
		return fmt.Errorf("create records_fetched_total counter: %w", err)
	}
	if t.recordsPostedTotal, err = t.meter.Int64Counter("erp_job_records_posted_total"); err != nil {
		return fmt.Errorf("create records_posted_total counter: %w", err)
	}
	if t.retriesTotal, err = t.meter.Int64Counter("erp_job_retries_total"); err != nil {
		return fmt.Errorf("create retries_total counter: %w", err)
	}
	if t.failuresTotal, err = t.meter.Int64Counter("erp_job_failures_total"); err != nil {
		return fmt.Errorf("create failures_total counter: %w", err)
	}
	if t.checkpointLag, err = t.meter.Int64Histogram("erp_job_checkpoint_lag"); err != nil {
		return fmt.Errorf("create checkpoint_lag histogram: %w", err)
	}
	if t.lastSuccessGauge, err = t.meter.Int64ObservableGauge("erp_job_last_success_timestamp"); err != nil {
		return fmt.Errorf("create last_success_timestamp gauge: %w", err)
	}

	if _, err := t.meter.RegisterCallback(func(_ context.Context, observer metric.Observer) error {
		value := t.lastSuccessUnix.Load()
		if value > 0 {
			observer.ObserveInt64(t.lastSuccessGauge, value)
		}
		return nil
	}, t.lastSuccessGauge); err != nil {
		return fmt.Errorf("register telemetry callback: %w", err)
	}

	return nil
}

func (t *Telemetry) Tracer(name string) trace.Tracer {
	if name == "" {
		return t.tracer
	}
	return otel.Tracer(name)
}

func (t *Telemetry) StartRun(ctx context.Context, runID string) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, "erp-job.run", trace.WithAttributes(attribute.String("run.id", runID)))
}

func (t *Telemetry) StartStep(ctx context.Context, step string) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, "erp-job.step", trace.WithAttributes(attribute.String("step", step)))
}

func (t *Telemetry) RecordRun(ctx context.Context, result string, duration time.Duration) {
	attrs := metric.WithAttributes(attribute.String("result", result))
	t.runsTotal.Add(ctx, 1, attrs)
	t.runDuration.Record(ctx, duration.Seconds(), attrs)
	if result == "success" {
		t.lastSuccessUnix.Store(time.Now().Unix())
	}
}

func (t *Telemetry) RecordFetched(ctx context.Context, step string, count int) {
	t.recordsFetchedTotal.Add(ctx, int64(count), metric.WithAttributes(attribute.String("step", step)))
}

func (t *Telemetry) RecordPosted(ctx context.Context, operation string, count int) {
	t.recordsPostedTotal.Add(ctx, int64(count), metric.WithAttributes(attribute.String("operation", operation)))
}

func (t *Telemetry) RecordRetry(ctx context.Context, endpointGroup string, operation string) {
	t.retriesTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String("endpoint_group", endpointGroup),
		attribute.String("operation", operation),
	))
}

func (t *Telemetry) RecordFailure(ctx context.Context, endpointGroup string, operation string, statusCode int, errorClass string) {
	attributes := []attribute.KeyValue{
		attribute.String("endpoint_group", endpointGroup),
		attribute.String("operation", operation),
		attribute.String("error_class", errorClass),
	}
	if statusCode > 0 {
		attributes = append(attributes, attribute.Int("status_code", statusCode))
	}
	t.failuresTotal.Add(ctx, 1, metric.WithAttributes(attributes...))
}

func (t *Telemetry) RecordCheckpointLag(ctx context.Context, entity string, lag int) {
	t.checkpointLag.Record(ctx, int64(lag), metric.WithAttributes(attribute.String("entity", entity)))
}

func WithAttemptObserver(ctx context.Context, observer AttemptObserver) context.Context {
	return context.WithValue(ctx, attemptObserverKey{}, observer)
}

func AttemptObserverFromContext(ctx context.Context) AttemptObserver {
	observer, _ := ctx.Value(attemptObserverKey{}).(AttemptObserver)
	return observer
}

func WithRunID(ctx context.Context, runID string) context.Context {
	return context.WithValue(ctx, runIDKey{}, runID)
}

func RunIDFromContext(ctx context.Context) string {
	runID, _ := ctx.Value(runIDKey{}).(string)
	return runID
}

func NewRunID() (string, error) {
	var buffer [16]byte
	if _, err := rand.Read(buffer[:]); err != nil {
		return "", fmt.Errorf("generate run id: %w", err)
	}
	return hex.EncodeToString(buffer[:]), nil
}

func exporterOptions(cfg config.OTel) ([]otlptracehttp.Option, []otlpmetrichttp.Option, error) {
	parsed, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("parse OTEL endpoint: %w", err)
	}

	endpoint := parsed.Host
	if endpoint == "" {
		endpoint = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(cfg.Endpoint, "https://"), "http://"))
	}

	traceOpts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(endpoint)}
	metricOpts := []otlpmetrichttp.Option{otlpmetrichttp.WithEndpoint(endpoint)}

	if parsed.Scheme == "http" || cfg.Insecure {
		traceOpts = append(traceOpts, otlptracehttp.WithInsecure())
		metricOpts = append(metricOpts, otlpmetrichttp.WithInsecure())
	}

	return traceOpts, metricOpts, nil
}
