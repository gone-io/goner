package grpc

import (
	"context"
	"strings"
	"testing"

	mock "github.com/gone-io/gone/mock/v2"
	"github.com/gone-io/gone/v2"
	gMock "github.com/gone-io/goner/g/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func (s *clientRegister) Infof(format string, args ...any) {}
func TestClientRegister_traceInterceptor(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()
	tracer := gMock.NewMockTracer(ctr)
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

func TestClientRegister_CreateConn(t *testing.T) {
	tests := []struct {
		name                string
		insecure            bool
		hasDiscovery        bool
		loadBalancingPolicy string
		expectedError       bool
	}{
		{
			name:          "基本连接配置",
			insecure:      true,
			expectedError: false,
		},
		{
			name:                "带服务发现和负载均衡的配置",
			insecure:            true,
			hasDiscovery:        true,
			loadBalancingPolicy: "round_robin",
			expectedError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			register := &clientRegister{
				insecure:            tt.insecure,
				loadBalancingPolicy: tt.loadBalancingPolicy,
			}

			if tt.hasDiscovery {
				controller := gomock.NewController(t)
				defer controller.Finish()
				discovery := gMock.NewMockServiceDiscovery(controller)
				register.discovery = discovery
				register.Init()
			}

			conn, err := register.createConn(":0")
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, conn)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, conn)
				conn.Close()
			}
		})
	}
}

func TestClientRegister_GetConn(t *testing.T) {
	tests := []struct {
		name          string
		address       string
		expectedError bool
	}{
		{
			name:          "获取已缓存的连接",
			address:       ":0",
			expectedError: false,
		},
		{
			name:          "创建新连接",
			address:       ":1",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			register := &clientRegister{
				insecure:    true,
				connections: make(map[string]*grpc.ClientConn),
			}

			// 第一次获取连接
			conn1, err := register.getConn(tt.address)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, conn1)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, conn1)

			// 第二次获取相同地址的连接
			conn2, err := register.getConn(tt.address)
			assert.NoError(t, err)
			assert.Equal(t, conn1, conn2, "应返回缓存的连接实例")

			// 清理
			conn1.Close()
		})
	}
}

func TestClientRegister_Provide(t *testing.T) {

	tests := []struct {
		name          string
		tagConf       string
		configValue   string
		expectedError bool
	}{
		{
			name:          "set address directly",
			tagConf:       `address=127.0.0.1:4451`,
			expectedError: false,
		},
		{
			name:          "read address from config",
			tagConf:       `config=grpc.address`,
			configValue:   ":1",
			expectedError: false,
		},
		{
			name:          "empty address",
			tagConf:       ``,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			configure := mock.NewMockConfigure(controller)
			if tt.configValue != "" {
				configure.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(key string, value *string, defaultValue string) error {
						*value = tt.configValue
						return nil
					})
			}

			register := &clientRegister{
				insecure:    true,
				connections: make(map[string]*grpc.ClientConn),
				configure:   configure,
			}

			conn, err := register.Provide(tt.tagConf)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, conn)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, conn)
				conn.Close()
			}
		})
	}
}
