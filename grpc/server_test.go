package gone_grpc

import (
	"context"
	"errors"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/tracer"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"testing"
)

func (s *server) Errorf(format string, args ...any) {}
func (s *server) Warnf(format string, args ...any)  {}
func (s *server) Infof(format string, args ...any)  {}
func (s *server) Go(fn func())                      {}

func Test_createListener(t *testing.T) {
	err := createListener(&server{})
	assert.Nil(t, err)
}

func Test_server_initListener(t *testing.T) {
	t.Run("use cMuxServer", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		cMuxServer := NewMockCMuxServer(controller)
		listener := NewMockListener(controller)
		cMuxServer.EXPECT().MatchWithWriters(gomock.Any()).Return(listener)
		cMuxServer.EXPECT().GetAddress().Return("")

		s := server{
			cMuxServer: cMuxServer,
		}
		err := s.initListener()
		assert.Nil(t, err)
		assert.NotNil(t, s.listener)
	})

	t.Run("use tcpListener", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		listener := NewMockListener(controller)

		s := server{
			createListener: func(s *server) error {
				s.listener = listener
				return nil
			},
		}
		err := s.initListener()
		assert.Nil(t, err)
		assert.NotNil(t, s.listener)
	})

	t.Run("use tcpListener error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		listener := NewMockListener(controller)

		s := server{
			createListener: func(s *server) error {
				s.listener = listener
				return errors.New("error")
			},
		}
		err := s.initListener()
		assert.Error(t, err)
	})
}

type addr struct{}

func (a *addr) Network() string {
	return "tcp"
}
func (a *addr) String() string {
	return ":8080"
}

func Test_server_server(t *testing.T) {
	t.Run("server", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		listener := NewMockListener(controller)
		listener.EXPECT().Addr().Return(&addr{}).AnyTimes()
		listener.EXPECT().Accept().Return(nil, errors.New("error"))
		listener.EXPECT().Close().Return(nil)

		gone.
			NewApp().
			Test(func(logger gone.Logger) {
				s := server{
					grpcServer: grpc.NewServer(),
					listener:   listener,
					logger:     logger,
				}
				s.server()
			})
	})
}

func Test_server_Stop(t *testing.T) {
	s := server{
		grpcServer: grpc.NewServer(),
	}
	err := s.Stop()
	assert.Nil(t, err)
}

func Test_server_traceInterceptor(t *testing.T) {
	ctx := context.Background()
	traceId := "trace"

	gone.
		NewApp(tracer.Load).
		Test(func(in struct {
			tracer      tracer.Tracer `gone:"*"`
			tracerIdKey string        `gone:"config,server.grpc.x-trace-id-key=X-Trace-Id"`
		}) {
			s := server{
				tracer:      in.tracer,
				tracerIdKey: in.tracerIdKey,
			}

			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				in.tracerIdKey: []string{traceId},
			})

			var req any
			_, err := s.traceInterceptor(ctx, req, nil, func(ctx context.Context, req any) (any, error) {
				id := in.tracer.GetTraceId()
				assert.Equal(t, traceId, id)
				return nil, nil
			})
			assert.Nil(t, err)
		})
}

func Test_server_recoveryInterceptor(t *testing.T) {
	gone.
		NewApp(tracer.Load).
		Test(func(tracer tracer.Tracer, logger gone.Logger) {
			s := server{
				tracer: tracer,
				logger: logger,
			}
			_, err := s.recoveryInterceptor(context.Background(), 1, nil,
				func(ctx context.Context, req any) (any, error) {
					if req == 1 {
						panic(errors.New("error"))
					}
					return nil, nil
				})
			assert.Error(t, err)
		})
}
