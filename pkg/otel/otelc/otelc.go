package otelc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Provider struct {
	grpcClient     *grpc.ClientConn
	res            *resource.Resource
	tracerProvider *trace.TracerProvider
	meterProvider  *metric.MeterProvider
}

func New(target string, resource *resource.Resource, opts ...grpc.DialOption) (*Provider, error) {
	p := &Provider{
		grpcClient:     nil,
		res:            resource,
		tracerProvider: nil,
		meterProvider:  nil,
	}
	if err := p.initConn(target, opts...); err != nil {
		return nil, err
	}

	if err := p.setup(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Provider) Shutdown(ctx context.Context) (errs error) {
	if p.meterProvider != nil {
		errors.Join(errs, p.meterProvider.Shutdown(ctx))
	}

	if p.tracerProvider != nil {
		errors.Join(errs, p.tracerProvider.Shutdown(ctx))
	}

	return
}

func (p *Provider) setup() error {
	ctx := context.Background()

	var err error
	if err = p.setupTracerProvider(ctx, p.res); err != nil {
		return err
	}

	if err = p.setupMeterProvider(ctx, p.res); err != nil {
		return err
	}

	return nil
}

func (p *Provider) setupTracerProvider(ctx context.Context, res *resource.Resource) error {
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(p.grpcClient))
	if err != nil {
		return fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := trace.NewBatchSpanProcessor(traceExporter)
	p.tracerProvider = trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(p.tracerProvider)

	// 将全局传播程序设置为tracecontext (默认为no-op)。
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return nil
}

func (p *Provider) setupMeterProvider(ctx context.Context, res *resource.Resource) error {
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(p.grpcClient))
	if err != nil {
		return fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	p.meterProvider = metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(p.meterProvider)

	return nil
}

func (p *Provider) initConn(target string, opts ...grpc.DialOption) (err error) {
	p.grpcClient, err = grpc.NewClient(target, opts...)
	if err != nil {
		return fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}
	// 连接
	p.grpcClient.Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("connect timeout after state")
		default:
			state := p.grpcClient.GetState()
			if state == connectivity.Ready {
				return nil
			}

			if !p.grpcClient.WaitForStateChange(ctx, state) {
				return fmt.Errorf("connect error state: %s", p.grpcClient.GetState())
			}
		}
	}
}
