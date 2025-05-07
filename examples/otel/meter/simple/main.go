package main

import (
	"context"
	"github.com/gone-io/gone/v2"
	"go.opentelemetry.io/otel/metric"
	"os"
	"time"
)

func main() {
	_ = os.Setenv("GONE_OTEL_SERVICE_NAME", "simple meter demo")

	gone.
		NewApp(GoneModuleLoad).
		Run(func(meter metric.Meter, logger gone.Logger) {
			apiCounter, err := meter.Int64Counter(
				"api.counter",
				metric.WithDescription("API调用的次数"),
				metric.WithUnit("{次}"),
			)
			if err != nil {
				logger.Errorf("create meter err:%v", err)
				return
			}

			for i := 0; i < 5; i++ {
				time.Sleep(1 * time.Second)
				apiCounter.Add(context.Background(), 1)
			}
		})
}
