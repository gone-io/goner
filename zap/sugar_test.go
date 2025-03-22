package gone_zap

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"os"
	"testing"
)

type tracer struct {
	gone.Flag
	g.Tracer
}

func (t *tracer) SetTraceId(traceId string, fn func()) {
	fn()
}

// GetTraceId Get the traceId of the current goroutine
func (t *tracer) GetTraceId() string {
	return "trace-id"
}

// Go Start a new goroutine instead of `go func`, which can pass the traceId to the new goroutine.
func (t *tracer) Go(fn func()) {
	go fn()
}

func TestNewSugar(t *testing.T) {
	os.Setenv("GONE_LOG_LEVEL", "debug")
	gone.
		NewApp(Load).
		Load(&tracer{}).
		Test(func(log gone.Logger, tracer g.Tracer, in struct {
			level string `gone:"config,log.level"`
		}) {
			log.Infof("level:%s", in.level)
			if in.level != "debug" {
				t.Fatal("log level error")
			}
			tracer.SetTraceId("", func() {
				log.Debugf("debug log")
				log.Infof("info log")
				log.Warnf("warn log")
				log.Errorf("error log")
			})
		})
}
