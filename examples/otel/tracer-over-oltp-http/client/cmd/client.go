package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/urllib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const tracerName = "hello-client"

//go:generate gonectr generate -m . -s .. -e client
func main() {
	gone.
		Load(&client{}).
		Run(func(c *client) {
			tracer := otel.Tracer(tracerName)
			ctx, span := tracer.Start(context.Background(), "RUN DEMO")
			defer func() {
				span.End()
			}()

			span.AddEvent("call server")
			c.CallServer(ctx)
		})
}

type client struct {
	gone.Flag
	client urllib.Client `gone:"*"`
	logger gone.Logger   `gone:"*"`
}

func (s *client) CallServer(ctx context.Context) {
	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, "CALL HELLO SERVER")
	defer span.End()

	var result urllib.Res[string]

	res, err := s.client.R().
		SetBody(map[string]any{
			"name": "jim",
		}).
		SetSuccessResult(&result).
		SetErrorResult(&result).
		SetContext(ctx).
		Post("http://127.0.0.1:8080/hello")

	if err != nil {
		s.logger.Errorf("client request err: %v", err)
		span.SetStatus(codes.Error, "call server failed")
		span.RecordError(err)
		return
	}
	s.logger.Infof("res.httpStatus=>%s", res.Status)
	s.logger.Infof("result=> %#v", result)
}
