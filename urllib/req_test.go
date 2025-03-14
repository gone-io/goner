package urllib

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/tracer"
	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_r_trip(t *testing.T) {
	gone.
		NewApp(tracer.Load).
		Test(func(in struct {
			tracer      tracer.Tracer `gone:"*"`
			tracerIdKey string        `gone:"config,urllib.req.x-trace-id-key=X-Trace-Id"`
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
