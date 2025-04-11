package openai

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var clientMap sync.Map

type Config struct {
	ApiToken         string         `json:"apiToken"`
	BaseUrl          string         `json:"baseUrl"`
	OrgID            string         `json:"orgID"`
	APIType          openai.APIType `json:"APIType"`
	APIVersion       string         `json:"APIVersion"`
	AssistantVersion string         `json:"assistantVersion"`
	ProxyUrl         string         `json:"proxyUrl"`
}

func (config Config) ToOpenAiConfig(httpDoer g.HTTPDoer) openai.ClientConfig {
	var conf openai.ClientConfig

	switch config.APIType {
	case openai.APITypeAzure:
		conf = openai.DefaultAzureConfig(config.ApiToken, config.BaseUrl)
	case openai.APITypeAnthropic:
		conf = openai.DefaultAnthropicConfig(config.ApiToken, config.BaseUrl)
	default:
		conf = openai.DefaultConfig(config.ApiToken)
	}

	if config.BaseUrl != "" {
		conf.BaseURL = config.BaseUrl
	}

	if config.ProxyUrl != "" {
		Url, err := url.Parse(config.ProxyUrl)
		if err != nil {
			panic(gone.ToError(err))
		}
		transport := http.Transport{
			Proxy: http.ProxyURL(Url),
		}
		conf.HTTPClient = &http.Client{
			Transport: &transport,
		}
	} else if httpDoer != nil {
		conf.HTTPClient = httpDoer
	}
	if config.APIType != "" {
		conf.APIType = config.APIType
	}
	if config.APIVersion != "" {
		conf.APIVersion = config.APIVersion
	}
	if config.AssistantVersion != "" {
		conf.AssistantVersion = config.AssistantVersion
	}
	return conf
}

var (
	provider = gone.WrapFunctionProvider(func(tagConf string, param struct {
		configure gone.Configure `gone:"configure"`
		httpDoer  g.HTTPDoer     `gone:"openai-proxies" option:"allowNil"`
	}) (*openai.Client, error) {
		if value, ok := clientMap.Load(tagConf); ok {
			return value.(*openai.Client), nil
		}

		var prefix = "openai"
		_, keys := gone.TagStringParse(tagConf)
		if len(keys) > 0 && keys[0] != "" {
			prefix = strings.TrimSpace(keys[0])
		}
		var config Config
		err := param.configure.Get(prefix, &config, "")
		if err != nil {
			return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("get %s config err", prefix))
		}

		client := openai.NewClientWithConfig(config.ToOpenAiConfig(param.httpDoer))
		clientMap.Store(tagConf, client)
		return client, nil
	})

	load = g.BuildOnceLoadFunc(g.L(provider))
)

func Load(loader gone.Loader) error {
	return load(loader)
}
