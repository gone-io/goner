package apollo

import (
	"github.com/apolloconfig/agollo/v4/agcache"
	"github.com/apolloconfig/agollo/v4/storage"
	"go.uber.org/mock/gomock"
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

func Test_apolloConfigure_Notify(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var l storage.ChangeListener

	mockClient := NewMockClient(controller)
	mockClient.EXPECT().AddChangeListener(gomock.Any()).DoAndReturn(func(listener storage.ChangeListener) {
		l = listener
	})

	applicationCache := NewMockCacheInterface(controller)
	applicationCache.EXPECT().Range(gomock.Any()).DoAndReturn(func(fn func(key, value any) bool) {
		fn("test.key", "test.value1")
	})

	testCache := NewMockCacheInterface(controller)
	testCache.EXPECT().Range(gomock.Any()).DoAndReturn(func(fn func(key, value any) bool) {
		fn("test.key", "test.value2")
	})

	mockClient.EXPECT().GetConfigCache(gomock.Any()).DoAndReturn(func(key string) agcache.CacheInterface {
		if key == "application" {
			return applicationCache
		}
		return testCache
	}).Times(2)

	startWithConfig = func(loadAppConfig func() (*config.AppConfig, error)) (agollo.Client, error) {
		return mockClient, nil
	}

	gone.
		NewApp(Load).
		Test(func(watcher gone.ConfWatcher) {
			key := "test.key"

			var oldVal, newVal any

			watcher(key, func(o, n any) {
				oldVal, newVal = o, n
			})

			event := storage.ChangeEvent{
				Changes: map[string]*storage.ConfigChange{
					key: {
						OldValue:   "test.value2",
						NewValue:   "test.updated",
						ChangeType: storage.MODIFIED,
					},
				},
			}
			event.Namespace = "test.yml"

			l.OnChange(&event)
			assert.Equal(t, "test.value2", oldVal)
			assert.Equal(t, "test.updated", newVal)

			event.Changes = map[string]*storage.ConfigChange{
				key: {
					OldValue:   "test.value1",
					NewValue:   "test.updated2",
					ChangeType: storage.MODIFIED,
				},
			}

			l.OnChange(&event)
			assert.Equal(t, "test.updated", oldVal)
			assert.Equal(t, "test.updated2", newVal)

			event.Changes = map[string]*storage.ConfigChange{
				key: {
					OldValue:   "test.value1",
					NewValue:   "test.updated2",
					ChangeType: storage.DELETED,
				},
			}

			l.OnChange(&event)
			assert.Equal(t, "test.updated2", oldVal)
			assert.Equal(t, "test.value1", newVal)

			watcher("test.key2", func(o, n any) {
				oldVal, newVal = o, n
			})

			event.Changes = map[string]*storage.ConfigChange{
				"test.key2": {
					OldValue:   "test.value2",
					NewValue:   "test.updated2",
					ChangeType: storage.ADDED,
				},
			}

			l.OnChange(&event)
			assert.Equal(t, nil, oldVal)
			assert.Equal(t, "test.updated2", newVal)

			event.Changes = map[string]*storage.ConfigChange{
				"test.key2": {
					OldValue:   "test.value2",
					NewValue:   "test.updated2",
					ChangeType: storage.DELETED,
				},
			}

			l.OnChange(&event)
			assert.Equal(t, "test.updated2", oldVal)
			assert.Equal(t, nil, newVal)
		})
}
