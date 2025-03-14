package gone_grpc

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/tracer"
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
	gone.
		NewApp(tracer.Load).
		Test(func(in struct {
			tracer      tracer.Tracer `gone:"gone-tracer"`
			tracerIdKey string        `gone:"config,server.grpc.x-trace-id-key=X-Trace-Id"`
		}) {

			var req, reply any

			register := clientRegister{
				tracer:      in.tracer,
				tracerIdKey: in.tracerIdKey,
			}

			tracer.SetTraceId("xxxx", func() {
				err := register.traceInterceptor(
					context.Background(),
					"test",
					req, reply,
					nil,
					func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
						md, b := metadata.FromOutgoingContext(ctx)
						assert.True(t, b)
						list := md[strings.ToLower(in.tracerIdKey)]

						assert.Equal(t, 1, len(list))

						assert.Equal(t, "xxxx", list[0])
						return nil
					},
				)
				assert.Nil(t, err)
			})
		})
}

func Test_clientRegister_register(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	client := NewMockClient(controller)
	client.EXPECT().Address().Return(":8080").AnyTimes()
	client.EXPECT().Stub(gomock.Any())

	register := clientRegister{
		connections: make(map[string]*grpc.ClientConn),
		clients:     []Client{client},
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
