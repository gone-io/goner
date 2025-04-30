package grpc

import (
	"context"
	"github.com/gone-io/goner/g"
	gMock "github.com/gone-io/goner/g/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestWithOpenTelemetry(t *testing.T) {
	s := server{
		createListener:     mustCreateListener,
		getLocalIps:        g.GetLocalIps,
		isOtelTracerLoaded: true,
	}
	s.Init()

	register := clientRegister{
		connections:        make(map[string]*grpc.ClientConn),
		isOtelTracerLoaded: true,
		insecure:           true,
	}

	conn, err := register.createConn("127.0.0.1:9090")
	assert.Nil(t, err)
	assert.NotNil(t, conn)
}

func TestWithTracer(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	tracer := gMock.NewMockTracer(controller)

	register := clientRegister{
		connections: make(map[string]*grpc.ClientConn),
		insecure:    true,
		tracer:      tracer,
		tracerIdKey: "X-Trace-Id",
	}
	traceId := "xxx-0001"
	tracer.EXPECT().GetTraceId().Return(traceId)
	err := register.traceInterceptor(context.Background(), "", nil, nil, nil, func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		incomingContext, b := metadata.FromOutgoingContext(ctx)
		assert.True(t, b)
		assert.Equal(t, traceId, incomingContext.Get(register.tracerIdKey)[0])
		return nil
	})
	assert.Nil(t, err)
}
