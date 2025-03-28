package gone_zap

import (
	"errors"
	"testing"

	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
)

// mockLoader 模拟Loader接口，用于测试加载失败的场景
type mockLoader struct {
	failOn string
}

func (m *mockLoader) Loaded(key gone.LoaderKey) bool {
	return false
}

func (m *mockLoader) Load(gone.Goner, ...gone.Option) error {
	if m.failOn == "atomicLevel" {
		return errors.New("failed to load atomicLevel")
	}
	if m.failOn == "zapLoggerProvider" {
		return errors.New("failed to load zapLoggerProvider")
	}
	if m.failOn == "sugarProvider" {
		return errors.New("failed to load sugarProvider")
	}
	if m.failOn == "sugar" {
		return errors.New("failed to load sugar")
	}
	return nil
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		failOn  string
		wantErr bool
	}{
		{"success case", "", false},
		{"fail on atomicLevel", "atomicLevel", true},
		{"fail on zapLoggerProvider", "zapLoggerProvider", true},
		{"fail on sugarProvider", "sugarProvider", true},
		{"fail on sugar", "sugar", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLoader := &mockLoader{failOn: tt.failOn}
			err := Load(mockLoader)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
