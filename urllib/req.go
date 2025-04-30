package urllib

import (
	"context"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/imroc/req/v3"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"net/http/httptrace"
	"path/filepath"
)

type r struct {
	gone.Flag
	*req.Client
	logger          gone.Logger          `gone:"*"`
	tracer          g.Tracer             `gone:"*" option:"allowNil"`
	lb              g.LoadBalancer       `gone:"*" option:"allowNil"`
	isOtelLogLoaded g.IsOtelTracerLoaded `gone:"*" option:"allowNil"`

	innerServicePattern string `gone:"config,urllib.inner-service-pattern=*"`
	requestIdKey        string `gone:"config,urllib.req.x-request-id-key=X-Request-Id"`
	tracerIdKey         string `gone:"config,urllib.req.x-trace-id-key=X-Trace-Id"`
}

var tracerName = "urllib"

func (r *r) trip(rt req.RoundTripper) req.RoundTripFunc {
	tracer := otel.Tracer(tracerName)

	return func(req *req.Request) (resp *req.Response, err error) {
		var ctx context.Context
		var span trace.Span
		if r.isOtelLogLoaded {
			ctx, span = tracer.Start(req.Context(), fmt.Sprintf("%s %s", req.Method, req.URL.Path))
			defer span.End()

			span.SetAttributes(
				attribute.String("http.url", req.URL.String()),
				attribute.String("http.method", req.Method),
				attribute.String("http.req.header", req.HeaderToString()),
			)
			if len(req.Body) > 0 {
				span.SetAttributes(
					attribute.String("http.req.body", string(req.Body)),
				)
			}
			req.SetContext(ctx)
		}

		matched, err := filepath.Match(r.innerServicePattern, req.URL.Host)
		if err != nil {
			r.logger.Errorf("match inner service err: %v", err.Error())
			return nil, gone.ToErrorWithMsg(err, "match inner service err")
		}

		if matched {
			if r.lb != nil {
				instance, err := r.lb.GetInstance(req.Context(), req.URL.Host)
				if err != nil {
					r.logger.Errorf("lb get instance err: %v", err)
					return nil, gone.ToError(err)
				}
				req.URL.Host = fmt.Sprintf("%s:%d", instance.GetIP(), instance.GetPort())
			}

			tracerId, _ := req.Context().Value(r.tracerIdKey).(string)
			if r.tracer != nil {
				tracerId = r.tracer.GetTraceId()
			}

			//传递traceId
			req.SetHeader(r.tracerIdKey, tracerId)
		}

		resp, err = rt.RoundTrip(req)

		if span != nil && r.isOtelLogLoaded {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			if resp.Response != nil {
				span.SetAttributes(
					attribute.Int("http.status_code", resp.StatusCode),
					attribute.String("http.resp.header", resp.HeaderToString()),
					attribute.String("http.resp.body", resp.String()),
				)
			}
		}
		return
	}
}

func (r *r) Init() error {
	r.Client = r.C()
	return nil
}

func (r *r) C() *req.Client {
	client := req.C()
	if r.isOtelLogLoaded {
		client.Transport.WrapRoundTripFunc(func(rt http.RoundTripper) req.HttpRoundTripFunc {
			return otelhttp.NewTransport(rt,
				otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
					return otelhttptrace.NewClientTrace(ctx)
				}),
			).RoundTrip
		})
	}
	client.WrapRoundTripFunc(r.trip)
	return client
}
