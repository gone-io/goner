package urllib

import (
	"errors"
	"net/url"
	"testing"

	"github.com/gone-io/gone/v2"
	gonemock "github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/gone-io/goner/g/mock"
	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_r_trip(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	mockTracer := mock.NewMockTracer(ctr)
	mockTracer.EXPECT().SetTraceId(gomock.Any(), gomock.Any()).DoAndReturn(func(traceId string, fn func()) {
		fn()
	}).AnyTimes()
	mockTracer.EXPECT().GetTraceId().AnyTimes().Return("xxxx")

	gone.
		NewApp(func(loader gone.Loader) error {
			return loader.Load(gone.WrapFunctionProvider(func(tagConf string, param struct{}) (g.Tracer, error) {
				return mockTracer, nil
			}))
		}).
		Test(func(in struct {
			tracer      g.Tracer `gone:"*"`
			tracerIdKey string   `gone:"config,urllib.req.x-trace-id-key=X-Trace-Id"`
		}) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tripper := NewMockRoundTripper(controller)

			in.tracer.SetTraceId("xxxx", func() {
				tripper.
					EXPECT().
					RoundTrip(gomock.Any()).
					Do(func(req *req.Request) {
						traceId := req.Headers.Get(in.tracerIdKey)
						assert.Equal(t, "xxxx", traceId)
					}).
					Return(nil, nil)

				g := r{
					tracer:              in.tracer,
					tracerIdKey:         in.tracerIdKey,
					innerServicePattern: "*.service",
				}
				trip := g.trip(tripper)
				_, err := trip(&req.Request{
					URL: &url.URL{
						Scheme: "https",
						Host:   "inner.service",
					},
				})
				assert.Nil(t, err)
			})
		})
}

func Test_r_trip_InnerService(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	service := mock.NewMockService(ctr)
	service.EXPECT().GetIP().Return("192.168.1.1")
	service.EXPECT().GetPort().Return(8080)

	balancer := mock.NewMockLoadBalancer(ctr)
	balancer.EXPECT().GetInstance(gomock.Any(), "internal.service").Return(service, nil)

	balancer.EXPECT().GetInstance(gomock.Any(), "internal.service").Return(service, errors.New("err"))

	tests := []struct {
		name                string
		innerServicePattern string
		url                 string
		lb                  g.LoadBalancer
		expectedHost        string
		expectError         bool
	}{
		{
			name:                "not inner service",
			innerServicePattern: "internal.*",
			url:                 "https://example.com",
			expectedHost:        "example.com",
			expectError:         false,
		},
		{
			name:                "inner service",
			innerServicePattern: "internal.*",
			url:                 "https://internal.service",
			lb:                  balancer,
			expectedHost:        "192.168.1.1:8080",
			expectError:         false,
		},
		{
			name:                "inner service get instance error",
			innerServicePattern: "internal.*",
			url:                 "https://internal.service",
			lb:                  balancer,
			expectedHost:        "192.168.1.1:8080",
			expectError:         true,
		},
		{
			name:                "bad service name",
			innerServicePattern: "[",
			url:                 "https://example.com",
			expectError:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tripper := NewMockRoundTripper(controller)
			tripper.EXPECT().RoundTrip(gomock.Any()).Return(nil, nil).AnyTimes()

			logger := goneMock.NewMockLogger(controller)
			logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

			// 初始化请求处理器
			r := &r{
				lb:                  tt.lb,
				logger:              logger,
				innerServicePattern: tt.innerServicePattern,
			}

			// 创建请求
			parsedURL, _ := url.Parse(tt.url)
			request := &req.Request{
				URL: parsedURL,
			}

			// 执行trip函数
			tripFunc := r.trip(tripper)
			_, err := tripFunc(request)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedHost, request.URL.Host)
		})
	}
}

func Test_r_trip_Error(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	mockTracer := mock.NewMockTracer(ctr)
	mockTracer.EXPECT().SetTraceId(gomock.Any(), gomock.Any()).AnyTimes()
	mockTracer.EXPECT().GetTraceId().AnyTimes().Return("xxxx")

	gone.
		NewApp(func(loader gone.Loader) error {
			return loader.Load(gone.WrapFunctionProvider(func(tagConf string, param struct{}) (g.Tracer, error) {
				return mockTracer, nil
			}))
		}).
		Test(func(in struct {
			tracer      g.Tracer `gone:"*"`
			tracerIdKey string   `gone:"config,urllib.req.x-trace-id-key=X-Trace-Id"`
		}) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tripper := NewMockRoundTripper(controller)

			in.tracer.SetTraceId("xxxx", func() {
				tripper.
					EXPECT().
					RoundTrip(gomock.Any()).
					Do(func(req *req.Request) {
						traceId := req.Headers.Get(in.tracerIdKey)
						assert.Equal(t, "xxxx", traceId)
					}).
					Return(nil, errors.New("round trip error"))

				g2 := r{
					tracer:      in.tracer,
					tracerIdKey: in.tracerIdKey,
				}
				trip := g2.trip(tripper)
				_, err := trip(&req.Request{})
				assert.Error(t, err)
			})
		})
}
