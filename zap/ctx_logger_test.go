package gone_zap

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	gMock "github.com/gone-io/goner/g/mock"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"testing"
)

func Test_ctxLogger_Ctx(t *testing.T) {

	provider := &zapLoggerProvider{
		zapLogger: zap.NewNop(),
	}

	traceId, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanContext := trace.ContextWithRemoteSpanContext(context.Background(), trace.SpanContext{}.WithTraceID(
		traceId,
	))

	controller := gomock.NewController(t)
	defer controller.Finish()
	gTracer := gMock.NewMockTracer(controller)

	type fields struct {
		tracer g.Tracer
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		before func()
		fields fields
		args   args
		want   func(logger gone.Logger)
	}{
		{
			name: "with span context",
			fields: fields{
				tracer: nil,
			},
			args: args{
				ctx: spanContext,
			},
			want: func(logger gone.Logger) {
				logger.Infof("traceId setted by ctx logger")
			},
		},
		{
			name: "without span context",
			fields: fields{
				tracer: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			want: func(logger gone.Logger) {
				logger.Infof("traceId setted by ctx logger")
			},
		},
		{
			name: "with tracer",
			before: func() {
				gTracer.EXPECT().GetTraceId().Return("4bf92f3577b34da6a3ce929d0e0e4736")
			},
			fields: fields{
				tracer: gTracer,
			},
			args: args{
				ctx: context.Background(),
			},
			want: func(logger gone.Logger) {
				logger.Infof("traceId setted by ctx logger")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before()
			}
			c := ctxLogger{
				provider: provider,
				tracer:   tt.fields.tracer,
			}
			tt.want(c.Ctx(tt.args.ctx))
		})
	}
}
