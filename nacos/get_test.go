package nacos

import (
	mock "github.com/gone-io/gone"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_configure_Get(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	mockConfigure := mock.NewMockConfigure(ctr)

	t.Run("Watch Mode", func(t *testing.T) {
		c := configure{
			watch:  true,
			keyMap: make(map[string][]any),
			viper:  viper.New(),
		}

		var value string
		c.viper.Set("test.key", "test-value")
		err := c.Get("test.key", &value, "")

		assert.NoError(t, err)
		assert.Equal(t, "test-value", value)
		assert.Contains(t, c.keyMap, "test.key")
		assert.Contains(t, c.keyMap["test.key"], &value)
	})

	t.Run("Empty Viper", func(t *testing.T) {
		c := configure{
			localConfigure: mockConfigure,
		}

		var value string
		mockConfigure.EXPECT().
			Get("test.key", &value, "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*string)) = "local-value"
				return nil
			})

		err := c.Get("test.key", &value, "")

		assert.NoError(t, err)
		assert.Equal(t, "local-value", value)
	})

	t.Run("Nacos Config", func(t *testing.T) {
		c := configure{
			viper: viper.New(),
		}

		var value string
		c.viper.Set("test.key", "nacos-value")
		err := c.Get("test.key", &value, "")

		assert.NoError(t, err)
		assert.Equal(t, "nacos-value", value)
	})

	t.Run("Empty Nacos Config", func(t *testing.T) {
		c := configure{
			viper:          viper.New(),
			localConfigure: mockConfigure,
		}

		var value string
		mockConfigure.EXPECT().
			Get("test.key", &value, "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*string)) = "local-value"
				return nil
			})

		err := c.Get("test.key", &value, "")

		assert.NoError(t, err)
		assert.Equal(t, "local-value", value)
	})

	t.Run("Unmarshal Error", func(t *testing.T) {
		c := configure{
			viper:                     viper.New(),
			localConfigure:            mockConfigure,
			useLocalConfIfKeyNotExist: true,
		}

		type Config struct {
			Value int
		}
		var value Config
		c.viper.Set("test.key", "invalid-value")

		mockConfigure.EXPECT().
			Get("test.key", gomock.Any(), "").
			DoAndReturn(func(key string, v any, defaultVal string) error {
				*(v.(*Config)) = Config{Value: 123}
				return nil
			})

		err := c.Get("test.key", &value, "")

		assert.NoError(t, err)
		assert.Equal(t, 123, value.Value)
	})
}
