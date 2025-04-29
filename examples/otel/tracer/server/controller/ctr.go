package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.opentelemetry.io/otel"
)

type ctr struct {
	gone.Flag
	r g.IRoutes `gone:"*"`
}

const tracerName = "hello-server"

func (c *ctr) Mount() (err g.MountError) {
	tracer := otel.Tracer(tracerName)

	c.r.POST("/hello", func(ctx *gin.Context, i struct {
		req struct {
			Name string `json:"name"`
		} `gone:"http,body"`
	}) string {
		x, span := tracer.Start(ctx, "hello")
		defer span.End()

		return SayHello(x, i.req.Name)
	})
	return
}

func SayHello(context context.Context, name string) string {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(context, "say-hello")
	defer span.End()
	return fmt.Sprintf("hello, %s", name)
}
