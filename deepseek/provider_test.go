package deepseek

import (
	"os"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"

	"github.com/stretchr/testify/assert"
)

type mockHTTPDoer struct {
	g.HTTPDoer
}

func TestConfig_ToDeepseekOptions(t *testing.T) {
	t.Run("test base url config", func(t *testing.T) {
		config := Config{
			AuthToken: "test-token",
			BaseURL:   "https://test.deepseek.com",
		}
		options := config.ToDeepseekOptions(nil)
		client, err := deepseek.NewClientWithOptions("test-token", options...)
		assert.NoError(t, err)
		assert.Equal(t, "https://test.deepseek.com", client.BaseURL)
	})

	t.Run("test timeout config", func(t *testing.T) {
		timeout := 30
		config := Config{
			AuthToken: "test-token",
			Timeout:   30,
		}
		options := config.ToDeepseekOptions(nil)
		client, err := deepseek.NewClientWithOptions("test-token", options...)
		assert.NoError(t, err)
		assert.Equal(t, timeout, int(client.Timeout))
	})

	t.Run("test path config", func(t *testing.T) {
		config := Config{
			AuthToken: "test-token",
			Path:      "custom/path",
		}
		options := config.ToDeepseekOptions(nil)
		client, err := deepseek.NewClientWithOptions("test-token", options...)
		assert.NoError(t, err)
		assert.Equal(t, "custom/path", client.Path)
	})

	t.Run("test default path", func(t *testing.T) {
		config := Config{
			AuthToken: "test-token",
		}
		options := config.ToDeepseekOptions(nil)
		client, err := deepseek.NewClientWithOptions("test-token", options...)
		assert.NoError(t, err)
		assert.Equal(t, "chat/completions", client.Path)
	})

	t.Run("test proxy config", func(t *testing.T) {
		config := Config{
			AuthToken: "test-token",
			ProxyUrl:  "http://proxy.test.com",
		}
		options := config.ToDeepseekOptions(nil)
		// Just verify options are created without error
		_, err := deepseek.NewClientWithOptions("test-token", options...)
		assert.NoError(t, err)
	})

	t.Run("test proxy config error", func(t *testing.T) {
		config := Config{
			AuthToken: "test-token",
			ProxyUrl:  "http://user^:passwo^rd@foo.com/",
		}
		defer func() {
			if err := recover(); err == nil {
				t.Error("error")
			}
		}()
		config.ToDeepseekOptions(nil)
	})

	t.Run("test ToDeepseekOptions with httpDoer", func(t *testing.T) {
		config := Config{
			AuthToken: "test-token",
		}
		doer := mockHTTPDoer{}
		options := config.ToDeepseekOptions(doer)
		// Just verify options are created without error
		_, err := deepseek.NewClientWithOptions("test-token", options...)
		assert.NoError(t, err)
	})
}

func TestLoad(t *testing.T) {
	t.Run("config error", func(t *testing.T) {
		_ = os.Setenv("GONE_DEEPSEEK", "--")
		defer func() {
			_ = os.Unsetenv("GONE_DEEPSEEK")
		}()

		defer func() {
			if err := recover(); err != nil {
				if err.(gone.Error).Code() != gone.NotSupport {
					t.Error(err)
				}
			}
		}()
		gone.NewApp(Load).Run(func(in struct {
			client *deepseek.Client `gone:"*"`
		}) {

		})
	})

	t.Run("load default", func(t *testing.T) {
		gone.NewApp(Load).Run(func(in struct {
			client *deepseek.Client `gone:"*"`
		}) {

		})
	})

	t.Run("load multi", func(t *testing.T) {
		gone.NewApp(Load).Run(func(in struct {
			client  *deepseek.Client `gone:"*"`
			client1 *deepseek.Client `gone:"*,baidu"`
			client2 *deepseek.Client `gone:"*,aliyun"`
		}) {

		})
	})
}
