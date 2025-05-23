package openai

import (
	"github.com/gone-io/goner/g"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gone-io/gone/v2"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func getOpenaiConfigToken(config openai.ClientConfig) string {
	of := reflect.ValueOf(config)
	name := of.FieldByName("authToken")
	return name.String()
}

type mockHttpDoer struct {
	g.HTTPDoer
}

func TestConfig_ToOpenAiConfig(t *testing.T) {
	t.Run("test azure config", func(t *testing.T) {
		config := Config{
			ApiToken: "test-token",
			BaseUrl:  "https://test.azure.com",
			APIType:  openai.APITypeAzure,
		}
		conf := config.ToOpenAiConfig(nil)
		assert.Equal(t, openai.APITypeAzure, conf.APIType)
		assert.Equal(t, "test-token", getOpenaiConfigToken(conf))
		assert.Equal(t, "https://test.azure.com", conf.BaseURL)
	})

	t.Run("test anthropic config", func(t *testing.T) {
		config := Config{
			ApiToken: "test-token",
			BaseUrl:  "https://test.anthropic.com",
			APIType:  openai.APITypeAnthropic,
		}
		conf := config.ToOpenAiConfig(nil)
		assert.Equal(t, openai.APITypeAnthropic, conf.APIType)
		assert.Equal(t, "test-token", getOpenaiConfigToken(conf))
		assert.Equal(t, "https://test.anthropic.com", conf.BaseURL)
	})

	t.Run("test default config", func(t *testing.T) {
		config := Config{
			ApiToken: "test-token",
		}
		conf := config.ToOpenAiConfig(nil)
		assert.Equal(t, "test-token", getOpenaiConfigToken(conf))
	})

	t.Run("test proxy config", func(t *testing.T) {
		config := Config{
			ApiToken: "test-token",
			ProxyUrl: "http://proxy.test.com",
		}
		conf := config.ToOpenAiConfig(nil)
		assert.NotNil(t, conf.HTTPClient)
		assert.NotNil(t, conf.HTTPClient.(*http.Client).Transport)
	})

	t.Run("test proxy config error", func(t *testing.T) {
		config := Config{
			ApiToken: "test-token",
			ProxyUrl: "http://user^:passwo^rd@foo.com/",
		}
		defer func() {
			if err := recover(); err == nil {
				t.Error("error")
			}
		}()
		config.ToOpenAiConfig(nil)
	})

	t.Run("test api version and assistant version", func(t *testing.T) {
		config := Config{
			ApiToken:         "test-token",
			APIVersion:       "v2",
			AssistantVersion: "assistant-v2",
		}
		conf := config.ToOpenAiConfig(nil)
		assert.Equal(t, "v2", conf.APIVersion)
		assert.Equal(t, "assistant-v2", conf.AssistantVersion)
	})
	t.Run("test ToOpenAiConfig with httpDoer", func(t *testing.T) {
		config := Config{
			ApiToken:         "test-token",
			APIVersion:       "v2",
			AssistantVersion: "assistant-v2",
		}
		doer := mockHttpDoer{}
		conf := config.ToOpenAiConfig(doer)
		assert.Equal(t, doer, conf.HTTPClient)
	})
}

func TestLoad(t *testing.T) {
	t.Run("config error", func(t *testing.T) {
		_ = os.Setenv("GONE_OPENAI", "--")
		defer func() {
			_ = os.Unsetenv("GONE_OPENAI")
		}()

		defer func() {
			if err := recover(); err != nil {
				if !strings.Contains(err.(gone.Error).Error(), "config err") {
					t.Error(err)
				}
			}
		}()
		gone.NewApp(Load).Run(func(in struct {
			client *openai.Client `gone:"*"`
		}) {

		})
	})

	t.Run("load default", func(t *testing.T) {
		gone.NewApp(Load).Run(func(in struct {
			client *openai.Client `gone:"*"`
		}) {

		})
	})
	t.Run("load multi", func(t *testing.T) {
		gone.NewApp(Load).Run(func(in struct {
			client  *openai.Client `gone:"*"`
			client1 *openai.Client `gone:"*,baidu"`
			client2 *openai.Client `gone:"*,aliyun"`
		}) {

		})
	})
}
