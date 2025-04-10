package g

import (
	"github.com/gone-io/gone/v2"
	"testing"
)

func TestGetLocalIps(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetLocalIps()
		})
	}
}

func TestRecover(t *testing.T) {
	type args struct {
		logger gone.Logger
	}
	tests := []struct {
		name string
		args args
		fn   func()
	}{
		{
			"test",
			args{
				logger: gone.GetDefaultLogger(),
			},
			func() {
				panic("test")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer Recover(tt.args.logger)
			tt.fn()
		})
	}
}
