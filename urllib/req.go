package urllib

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/tracer"
	"github.com/imroc/req/v3"
)

var load = gone.OnceLoad(func(loader gone.Loader) error {
	err := tracer.Load(loader)
	if err != nil {
		return gone.ToError(err)
	}
	return loader.Load(
		&r{},
		gone.IsDefault(new(Client)),
	)
})

func Load(loader gone.Loader) error {
	return load(loader)
}

// Priest Deprecated, use Load instead
func Priest(loader gone.Loader) error {
	return Load(loader)
}

type r struct {
	gone.Flag
	*req.Client
	tracer.Tracer `gone:"*"`

	requestIdKey string `gone:"config,urllib.req.x-request-id-key=X-Request-Id"`
	tracerIdKey  string `gone:"config,urllib.req.x-trace-id-key=X-Trace-Id"`
}

func (r *r) trip(rt req.RoundTripper) req.RoundTripFunc {
	return func(req *req.Request) (resp *req.Response, err error) {
		tracerId := r.GetTraceId()
		//传递traceId
		req.SetHeader(r.tracerIdKey, tracerId)
		resp, err = rt.RoundTrip(req)
		return
	}
}

func (r *r) Init() error {
	r.Client = req.C()
	r.Client.WrapRoundTripFunc(r.trip)
	return nil
}

func (r *r) C() *req.Client {
	c := req.C()
	c.WrapRoundTripFunc(r.trip)
	return c
}
