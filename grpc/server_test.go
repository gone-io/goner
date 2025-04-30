package grpc

import (
	"context"
	"github.com/gone-io/gone/mock/v2"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"testing"

	"github.com/gone-io/goner/g"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

type mockTracer struct {
	g.Tracer
	traceId string
	called  bool
}

func (m *mockTracer) SetTraceId(traceId string, f func()) {
	m.traceId = traceId
	m.called = true
	f()
}

func TestServer_traceInterceptor(t *testing.T) {
	tests := []struct {
		name    string
		tracer  g.Tracer
		ctx     context.Context
		wantErr bool
	}{
		{
			name:   "tracer is nil",
			tracer: nil,
			ctx:    context.Background(),
		},
		{
			name:   "tracer exists with traceId",
			tracer: &mockTracer{},
			ctx:    metadata.NewIncomingContext(context.Background(), metadata.Pairs("X-Trace-Id", "test-trace-id")),
		},
		{
			name:   "tracer exists without traceId",
			tracer: &mockTracer{},
			ctx:    context.Background(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				tracer:      tt.tracer,
				tracerIdKey: "X-Trace-Id",
			}

			_, err := s.traceInterceptor(tt.ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, nil
			})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.tracer != nil {
				if md, ok := metadata.FromIncomingContext(tt.ctx); ok && len(md.Get("X-Trace-Id")) > 0 {
					assert.Equal(t, "test-trace-id", tt.tracer.(*mockTracer).traceId)
					assert.True(t, tt.tracer.(*mockTracer).called)
				} else {
					assert.Empty(t, tt.tracer.(*mockTracer).traceId)
					assert.True(t, tt.tracer.(*mockTracer).called)
				}
			}
		})
	}
}

func TestServer_Provide(t *testing.T) {
	tests := []struct {
		name       string
		grpcServer *grpc.Server
		wantErr    bool
	}{
		{
			name:       "normal case",
			grpcServer: grpc.NewServer(),
			wantErr:    false,
		},
		{
			name:       "grpcServer is nil",
			grpcServer: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				grpcServer: tt.grpcServer,
			}

			got, err := s.Provide()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.grpcServer, got)
			}
		})
	}
}

func Test_server_recoveryInterceptor(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := mock.NewMockLogger(controller)
	s := server{
		logger: logger,
	}

	t.Run("panic", func(t *testing.T) {
		logger.EXPECT().Errorf(gomock.Any(), gomock.Any())
		resp, err := s.recoveryInterceptor(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			panic("panic")
		})
		assert.Nil(t, resp)
		assert.Error(t, err)
	})
	t.Run("normal", func(t *testing.T) {
		resp, err := s.recoveryInterceptor(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			return "test", nil
		})
		assert.Equal(t, "test", resp)
		assert.NoError(t, err)
	})
}
