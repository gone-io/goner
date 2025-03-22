package urllib

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_r_trip(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	mockTracer := NewMockTracer(ctr)
	mockTracer.EXPECT().SetTraceId(gomock.Any(), gomock.Any()).AnyTimes()
	mockTracer.EXPECT().GetTraceId().AnyTimes().Return("xxxx")

	gone.
		NewApp(func(loader gone.Loader) error {
			return loader.Load(mockTracer)
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
					tracer:      in.tracer,
					tracerIdKey: in.tracerIdKey,
				}
				trip := g.trip(tripper)
				_, err := trip(&req.Request{})
				assert.Nil(t, err)
			})
		})
}
