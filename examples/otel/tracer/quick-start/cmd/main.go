package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel"
)

const tracerName = "demo"

//go:generate gonectl generate -m . -s ..
func main() {
	gone.Run(func() {
		tracer := otel.Tracer(tracerName)
		ctx, span := tracer.Start(context.Background(), "run demo")
		defer span.End()
		span.AddEvent("x event")
		doSomething(ctx)
	})
}

func doSomething(ctx context.Context) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, "doSomething")
	defer span.End()
}
