package apollo

import (
	"testing"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/gone-io/gone/v2"
	"github.com/stretchr/testify/assert"
)

type mockConfigure struct {
	getFunc func(key string, value any, defaultVal string) error
}

func (m *mockConfigure) Get(key string, value any, defaultVal string) error {
	return m.getFunc(key, value, defaultVal)
}

func TestApolloConfigureInit(t *testing.T) {
	tests := []struct {
		name           string
		getFunc        func(key string, value any, defaultVal string) error
		expectedError  bool
		expectedConfig *config.AppConfig
	}{
		{
			name: "success case",
			getFunc: func(key string, value any, defaultVal string) error {
				switch key {
				case "apollo.appId":
					*(value.(*string)) = "testApp"
				case "apollo.cluster":
					*(value.(*string)) = "testCluster"
				case "apollo.ip":
					*(value.(*string)) = "localhost:8080"
				case "apollo.namespace":
					*(value.(*string)) = "application"
				case "apollo.secret":
					*(value.(*string)) = "testSecret"
				case "apollo.isBackupConfig":
					*(value.(*bool)) = true
				case "apollo.watch":
					*(value.(*bool)) = false
				case "apollo.useLocalConfIfKeyNotExist":
					*(value.(*bool)) = true
				}
				return nil
			},
			expectedError: false,
			expectedConfig: &config.AppConfig{
				AppID:          "testApp",
				Cluster:        "testCluster",
				IP:             "localhost:8080",
				NamespaceName:  "application",
				Secret:         "testSecret",
				IsBackupConfig: true,
			},
		},
		{
			name: "error case - missing appId",
			getFunc: func(key string, value any, defaultVal string) error {
				return gone.NewError(400, "config not found", 400)
			},
			expectedError:  true,
			expectedConfig: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configure := &apolloConfigure{}
			mockConf := &mockConfigure{getFunc: tt.getFunc}

			// 保存原始的startWithConfig函数
			originalStartWithConfig := startWithConfig
			defer func() { startWithConfig = originalStartWithConfig }()

			// 模拟startWithConfig函数
			startWithConfig = func(loadAppConfig func() (*config.AppConfig, error)) (agollo.Client, error) {
				config, err := loadAppConfig()
				if err != nil {
					return nil, err
				}
				assert.Equal(t, tt.expectedConfig, config)
				return nil, nil
			}

			config, err := configure.init(mockConf)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedConfig, config)
			}
		})
	}
}
