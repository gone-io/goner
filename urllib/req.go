package urllib

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/imroc/req/v3"
	"path/filepath"
)

var load = gone.OnceLoad(func(loader gone.Loader) error {
	err := loader.Load(&r{}, gone.IsDefault(new(Client)))
	if err != nil {
		return gone.ToError(err)
	}
	err = loader.Load(&requestProvider{})
	if err != nil {
		return gone.ToError(err)
	}
	return loader.Load(&clientProvider{})
})

func Load(loader gone.Loader) error {
	return load(loader)
}

type r struct {
	gone.Flag
	*req.Client
	tracer g.Tracer       `gone:"*" option:"allowNil"`
	lb     g.LoadBalancer `gone:"*" option:"allowNil"`
	logger gone.Logger    `gone:"*"`

	innerServicePattern string `gone:"config,urllib.inner-service-pattern=*"`
	requestIdKey        string `gone:"config,urllib.req.x-request-id-key=X-Request-Id"`
	tracerIdKey         string `gone:"config,urllib.req.x-trace-id-key=X-Trace-Id"`
}

func (r *r) trip(rt req.RoundTripper) req.RoundTripFunc {
	return func(req *req.Request) (resp *req.Response, err error) {
		matched, err := filepath.Match(r.innerServicePattern, req.URL.Host)
		if err != nil {
			r.logger.Errorf("match inner service err: %v", err.Error())
			return nil, gone.ToErrorWithMsg(err, "match inner service err")
		}

		if matched {
			if r.lb != nil {
				instance, err := r.lb.GetInstance(req.Context(), req.URL.Host)
				if err != nil {
					r.logger.Errorf("lb get instance err: %v", err)
				}
				req.URL.Host = fmt.Sprintf("%s:%d", instance.GetIP(), instance.GetPort())
			}

			tracerId, _ := req.Context().Value(r.tracerIdKey).(string)
			if r.tracer != nil {
				tracerId = r.tracer.GetTraceId()
			}

			//传递traceId
			req.SetHeader(r.tracerIdKey, tracerId)
		}

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
