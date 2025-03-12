package gone_zap

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/tracer"
	"os"
	"testing"
)

func TestNewSugar(t *testing.T) {

	os.Setenv("GONE_LOG_LEVEL", "debug")
	gone.
		NewApp(Priest).
		Test(func(log gone.Logger, tracer tracer.Tracer, in struct {
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
