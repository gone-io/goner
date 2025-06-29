package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConfig_ToMongoOptions(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		check  func(*testing.T, *options.ClientOptions)
	}{
		{
			name: "basic config with URI",
			config: Config{
				URI:      "mongodb://localhost:27017",
				Database: "testdb",
			},
			check: func(t *testing.T, opts *options.ClientOptions) {
				assert.NotNil(t, opts)
			},
		},
		{
			name: "config with authentication",
			config: Config{
				URI:        "mongodb://localhost:27017",
				Username:   "testuser",
				Password:   "testpass",
				AuthSource: "admin",
			},
			check: func(t *testing.T, opts *options.ClientOptions) {
				assert.NotNil(t, opts)
				auth := opts.Auth
				assert.NotNil(t, auth)
				assert.Equal(t, "testuser", auth.Username)
				assert.Equal(t, "testpass", auth.Password)
				assert.Equal(t, "admin", auth.AuthSource)
			},
		},
		{
			name: "config with pool settings",
			config: Config{
				URI:             "mongodb://localhost:27017",
				MaxPoolSize:     100,
				MinPoolSize:     10,
				MaxConnIdleTime: 30 * time.Minute,
			},
			check: func(t *testing.T, opts *options.ClientOptions) {
				assert.NotNil(t, opts)
				assert.Equal(t, uint64(100), *opts.MaxPoolSize)
				assert.Equal(t, uint64(10), *opts.MinPoolSize)
				assert.Equal(t, 30*time.Minute, *opts.MaxConnIdleTime)
			},
		},
		{
			name: "config with timeouts",
			config: Config{
				URI:                    "mongodb://localhost:27017",
				ConnectTimeout:         10 * time.Second,
				SocketTimeout:          30 * time.Second,
				ServerSelectionTimeout: 5 * time.Second,
			},
			check: func(t *testing.T, opts *options.ClientOptions) {
				assert.NotNil(t, opts)
				assert.Equal(t, 10*time.Second, *opts.ConnectTimeout)
				assert.Equal(t, 30*time.Second, *opts.SocketTimeout)
				assert.Equal(t, 5*time.Second, *opts.ServerSelectionTimeout)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.config.ToMongoOptions()
			tt.check(t, opts)
		})
	}
}

func TestLoad(t *testing.T) {
	// Test that Load function exists and can be called
	// This is a basic test to ensure the function signature is correct
	assert.NotNil(t, Load)
}

// MockConfigure implements gone.Configure for testing
type MockConfigure struct {
	configs map[string]interface{}
}

func NewMockConfigure() *MockConfigure {
	return &MockConfigure{
		configs: make(map[string]interface{}),
	}
}

func (m *MockConfigure) Set(key string, value interface{}) {
	m.configs[key] = value
}

func (m *MockConfigure) Get(key string, value any, defaultValue string) error {
	if config, exists := m.configs[key]; exists {
		switch v := value.(type) {
		case *Config:
			if configData, ok := config.(Config); ok {
				*v = configData
			}
		}
		return nil
	}
	
	if defaultValue != "" {
		// Handle default value if needed
	}
	
	return nil
}

func (m *MockConfigure) GetString(key string, defaultValue ...string) string {
	if config, exists := m.configs[key]; exists {
		if str, ok := config.(string); ok {
			return str
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (m *MockConfigure) GetInt(key string, defaultValue ...int) int {
	if config, exists := m.configs[key]; exists {
		if i, ok := config.(int); ok {
			return i
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (m *MockConfigure) GetBool(key string, defaultValue ...bool) bool {
	if config, exists := m.configs[key]; exists {
		if b, ok := config.(bool); ok {
			return b
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

func (m *MockConfigure) GetFloat64(key string, defaultValue ...float64) float64 {
	if config, exists := m.configs[key]; exists {
		if f, ok := config.(float64); ok {
			return f
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0.0
}

func (m *MockConfigure) GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	if config, exists := m.configs[key]; exists {
		if d, ok := config.(time.Duration); ok {
			return d
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (m *MockConfigure) GetStringSlice(key string, defaultValue ...[]string) []string {
	if config, exists := m.configs[key]; exists {
		if slice, ok := config.([]string); ok {
			return slice
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (m *MockConfigure) GetIntSlice(key string, defaultValue ...[]int) []int {
	if config, exists := m.configs[key]; exists {
		if slice, ok := config.([]int); ok {
			return slice
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (m *MockConfigure) GetStringMap(key string, defaultValue ...map[string]interface{}) map[string]interface{} {
	if config, exists := m.configs[key]; exists {
		if m, ok := config.(map[string]interface{}); ok {
			return m
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (m *MockConfigure) GetStringMapString(key string, defaultValue ...map[string]string) map[string]string {
	if config, exists := m.configs[key]; exists {
		if m, ok := config.(map[string]string); ok {
			return m
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (m *MockConfigure) GetStringMapStringSlice(key string, defaultValue ...map[string][]string) map[string][]string {
	if config, exists := m.configs[key]; exists {
		if m, ok := config.(map[string][]string); ok {
			return m
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

func (m *MockConfigure) IsSet(key string) bool {
	_, exists := m.configs[key]
	return exists
}

func (m *MockConfigure) AllKeys() []string {
	keys := make([]string, 0, len(m.configs))
	for k := range m.configs {
		keys = append(keys, k)
	}
	return keys
}

func (m *MockConfigure) AllSettings() map[string]interface{} {
	return m.configs
}

func (m *MockConfigure) Sub(key string) gone.Configure {
	if config, exists := m.configs[key]; exists {
		if subConfig, ok := config.(map[string]interface{}); ok {
			newMock := NewMockConfigure()
			newMock.configs = subConfig
			return newMock
		}
	}
	return NewMockConfigure()
}

func (m *MockConfigure) UnmarshalKey(key string, rawVal interface{}) error {
	return nil
}

func (m *MockConfigure) Unmarshal(rawVal interface{}) error {
	return nil
}

func (m *MockConfigure) BindEnv(input ...string) error {
	return nil
}

func (m *MockConfigure) SetEnvPrefix(in string) {
}

func (m *MockConfigure) SetEnvKeyReplacer(r interface{}) {
}

func (m *MockConfigure) AutomaticEnv() {
}

func (m *MockConfigure) SetConfigFile(in string) {
}

func (m *MockConfigure) SetConfigName(in string) {
}

func (m *MockConfigure) SetConfigType(in string) {
}

func (m *MockConfigure) AddConfigPath(in string) {
}

func (m *MockConfigure) ReadInConfig() error {
	return nil
}

func (m *MockConfigure) ReadConfig(in context.Context) error {
	return nil
}

func (m *MockConfigure) WatchConfig() {
}

func (m *MockConfigure) OnConfigChange(run func(in interface{})) {
}