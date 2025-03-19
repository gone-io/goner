package gone_zap

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"go.uber.org/zap/zapcore"
	"reflect"
	"testing"
	"time"
)

type mockConfigure struct {
	gone.Flag
	gone.Configure
	level []any
}

func (m *mockConfigure) Get(key string, v any, defaultVal string) error {
	if key == "log.level" {
		m.level = append(m.level, v)
		fmt.Printf("1.->m.level:%v\n", m.level)
		return gone.SetValue(reflect.ValueOf(v), v, defaultVal)
	}
	return nil
}

func (m *mockConfigure) setLevel(level string) {
	for _, l := range m.level {
		fmt.Printf("2.->m.level:%v\n", l)
		s := l.(*string)
		*s = level
	}
}

func TestLevelChangeByConfigure(t *testing.T) {
	config := &mockConfigure{}

	gone.
		NewApp().
		Load(config, gone.Name(gone.ConfigureName),
			gone.IsDefault(new(gone.Configure)),
			gone.ForceReplace(),
		).
		Load(&atomicLevel{}).
		Test(func(in *atomicLevel) {
			if in.Enabled(zapcore.DebugLevel) {
				t.Error("debug level should not be enabled")
			}
			if !in.Enabled(zapcore.InfoLevel) {
				t.Error("info level should be enabled")
			}

			l := in.Level()
			if l != zapcore.InfoLevel {
				t.Error("level should be info")
			}

			config.setLevel("warn")
			time.Sleep(1 * time.Millisecond)

			l = in.Level()
			if l != zapcore.WarnLevel {
				t.Error("level should be warn")
			}
		})
}
