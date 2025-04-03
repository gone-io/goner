package grpc

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"strings"
	"testing"
)

func (s *clientRegister) Infof(format string, args ...any) {}
func TestClientRegister_traceInterceptor(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()
	tracer := NewMockTracer(ctr)
	tracer.EXPECT().SetTraceId(gomock.Any(), gomock.Any()).AnyTimes()
	tracer.EXPECT().GetTraceId().Return("xxxx").AnyTimes()

	const tracerIdKey = "X-Trace-Id"

	register := clientRegister{
		tracer:      tracer,
		tracerIdKey: tracerIdKey,
	}
	var req, reply any

	tracer.SetTraceId("xxxx", func() {
		err := register.traceInterceptor(
			context.Background(),
			"test",
			req, reply,
			nil,
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				md, b := metadata.FromOutgoingContext(ctx)
				assert.True(t, b)
				list := md[strings.ToLower(tracerIdKey)]

				assert.Equal(t, 1, len(list))

				assert.Equal(t, "xxxx", list[0])
				return nil
			},
		)
		assert.Nil(t, err)
	})
}

func Test_clientRegister_register(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	client := NewMockClient(controller)
	client.EXPECT().Address().Return(":0").AnyTimes()
	client.EXPECT().Stub(gomock.Any())

	register := clientRegister{
		connections: make(map[string]*grpc.ClientConn),
		clients:     []Client{client},
		logger:      gone.GetDefaultLogger(),
		insecure:    true,
	}

	err := register.Start()
	assert.Nil(t, err)
}

func Test_clientRegister_Stop(t *testing.T) {
	register := clientRegister{
		connections: make(map[string]*grpc.ClientConn),
	}
	conn, err2 := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err2)
	register.connections[":8080"] = conn

	err := register.Stop()
	assert.Nil(t, err)
}
