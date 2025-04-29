package main

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

// 使用 net/http/httptrace 深入 net/http 内部的链路追踪示例

const (
	serviceName     = "httpclient-Demo"
	peerServiceName = "blog"
	jaegerEndpoint  = "127.0.0.1:4318"
	blogURL         = "https://liwenzhou.com"
)

// newJaegerTraceProvider 创建一个 Jaeger Trace Provider
func newJaegerTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	// 创建一个使用 HTTP 协议连接本机Jaeger的 Exporter
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(jaegerEndpoint),
		otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}
	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(serviceName)))
	if err != nil {
		return nil, err
	}
	// 创建 Provider
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // 采样
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(time.Second)),
	)
	return traceProvider, nil
}

// initTracer 初始化 Tracer
func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	tp, err := newJaegerTraceProvider(ctx)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)
	return tp, nil
}

func main() {
	ctx := context.Background()
	tp, err := initTracer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = tp.Shutdown(ctx)
	}()

	// 创建 http client，配置trace
	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx)
			}),
		),
	}

	// 创建tracer
	tr := otel.Tracer("http-client")

	// 开启 span，PeerService 指要连接的目标服务
	ctx, span := tr.Start(ctx, "GET BLOG", trace.WithAttributes(semconv.PeerService(peerServiceName)))
	defer span.End()

	// 构建请求
	//req, _ := http.NewRequestWithContext(ctx, http.MethodGet, blogURL, nil)
	//// 发送请求
	//res, _ := client.Do(req)
	//
	//body, _ := io.ReadAll(res.Body)
	//_ = res.Body.Close()
	//
	//fmt.Printf("Response Received: %s\n", body)
}
