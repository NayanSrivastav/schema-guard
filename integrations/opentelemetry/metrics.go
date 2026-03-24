package opentelemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// PipelineMetrics tracks deep structural metric dimensions mapping straight to Datadog/Grafana OpenTelemetry collectors
type PipelineMetrics struct {
	meter             metric.Meter
	schemaPassCounter metric.Int64Counter
	schemaFailCounter metric.Int64Counter
	retryLatencies    metric.Float64Histogram
	tokenCostGauge    metric.Float64UpDownCounter // Native USD cost tracking
}

// InitMetrics bootstraps OpenTelemetry telemetry instruments and bindings targeting the Go LLM router
func InitMetrics(meterProvider metric.MeterProvider) (*PipelineMetrics, error) {
	if meterProvider == nil { // default to global singleton if omitted
		meterProvider = otel.GetMeterProvider()
	}
	meter := meterProvider.Meter("github.com/schemaguard/schemaguard")

	passCounter, err := meter.Int64Counter("schemaguard.validation.pass",
		metric.WithDescription("Number of successfully validated LLM outputs"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init pass counter: %w", err)
	}

	failCounter, err := meter.Int64Counter("schemaguard.validation.fail",
		metric.WithDescription("Number of structural LLM validation engine rejections"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init fail counter: %w", err)
	}

	latencies, err := meter.Float64Histogram("schemaguard.retry.latency",
		metric.WithDescription("Generation mapping latencies for execution loops (ms)"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init latency histogram: %w", err)
	}

	costGauge, err := meter.Float64UpDownCounter("schemaguard.tokens.cost_usd",
		metric.WithDescription("Tracks dynamically cumulative estimated AI model burn rates in USD"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init cost gauge: %w", err)
	}

	return &PipelineMetrics{
		meter:             meter,
		schemaPassCounter: passCounter,
		schemaFailCounter: failCounter,
		retryLatencies:    latencies,
		tokenCostGauge:    costGauge,
	}, nil
}

// RecordValidation fires metric mapping events capturing pure schema dimension flags
func (pm *PipelineMetrics) RecordValidation(ctx context.Context, passed bool) {
	if passed {
		pm.schemaPassCounter.Add(ctx, 1) // In production, attach attribute mappings (schema_name, etc)
	} else {
		pm.schemaFailCounter.Add(ctx, 1)
	}
}

// RecordLatency tracks exact sub-millisecond iteration timings
func (pm *PipelineMetrics) RecordLatency(ctx context.Context, ms float64) {
	pm.retryLatencies.Record(ctx, ms)
}

// RecordCost captures active token conversion burn values directly inside OpenTelemetry streams
func (pm *PipelineMetrics) RecordCost(ctx context.Context, usd float64) {
	pm.tokenCostGauge.Add(ctx, usd)
}
