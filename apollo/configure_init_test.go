package apollo

import (
	"testing"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/gone-io/gone/v2"
	originViper "github.com/spf13/viper"
	"go.uber.org/mock/gomock"
)

func TestApolloConfigureInitMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock dependencies
	mockClient := NewMockClient(ctrl)
	mockCache := NewMockCacheInterface(ctrl)
	mockLocalConfigure := NewMockConfigure(ctrl)
	mockViper := originViper.New()

	// Mock behavior
	oldStartWithConfig := startWithConfig
	oldViperNew := viperNew
	oldOriginViperNew := originViperNew
	defer func() {
		startWithConfig = oldStartWithConfig
		viperNew = oldViperNew
		originViperNew = oldOriginViperNew
	}()

	startWithConfig = func(configFunc func() (*config.AppConfig, error)) (agollo.Client, error) {
		return mockClient, nil
	}
	viperNew = func(testFlag gone.TestFlag) gone.Configure {
		return mockLocalConfigure
	}
	originViperNew = func() *originViper.Viper {
		return mockViper
	}

	watch := false

	mockLocalConfigure.EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(key string, value any, defaultVal string) {
			if key == "apollo.watch" {
				*value.(*bool) = watch
			}
		}).
		AnyTimes()

	// Test cases
	t.Run("success case with watch enabled", func(t *testing.T) {
		watch = true

		listener := changeListener{}
		listener.Init()

		configure := &apolloConfigure{
			testFlag:       nil,
			changeListener: &listener,
		}

		mockClient.EXPECT().GetConfigCache(gomock.Any()).Return(mockCache).AnyTimes()
		mockCache.EXPECT().Range(gomock.Any()).DoAndReturn(func(f func(key, value interface{}) bool) {
			f("key1", "value1")
			f("key2", "value2")
		}).AnyTimes()
		mockClient.EXPECT().AddChangeListener(gomock.Any()).Times(1)

		err := configure.Init()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("success case without watch", func(t *testing.T) {
		watch = false

		configure := &apolloConfigure{
			testFlag:       nil,
			changeListener: &changeListener{},
			watch:          false,
		}

		mockClient.EXPECT().GetConfigCache(gomock.Any()).Return(mockCache).AnyTimes()
		mockCache.EXPECT().Range(gomock.Any()).DoAndReturn(func(f func(key, value interface{}) bool) {
			f("key1", "value1")
			f("key2", "value2")
		}).AnyTimes()

		err := configure.Init()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("error case - client initialization fails", func(t *testing.T) {
		mockLocalConfigure.EXPECT().
			Get(gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		startWithConfig = func(configFunc func() (*config.AppConfig, error)) (agollo.Client, error) {
			return nil, gone.ToError("client init error")
		}

		configure := &apolloConfigure{
			testFlag:       nil,
			changeListener: &changeListener{},
		}

		err := configure.Init()
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
