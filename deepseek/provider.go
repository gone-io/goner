package deepseek

import (
	"fmt"
	"github.com/cohesion-org/deepseek-go"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Config struct {
	AuthToken string        `json:"authToken"` // The authentication token for the API
	BaseURL   string        `json:"baseURL"`   // The base URL for the API
	Timeout   time.Duration `json:"timeout"`   // The timeout for the current Client
	Path      string        `json:"path"`      // The path for the API request. Defaults to "chat/completions"
	ProxyUrl  string        `json:"proxyUrl"`
}

func (c Config) ToDeepseekOptions(httpDoer g.HTTPDoer) []deepseek.Option {
	options := make([]deepseek.Option, 0)
	if c.BaseURL != "" {
		options = append(options, deepseek.WithBaseURL(c.BaseURL))
	}
	if c.Timeout != 0 {
		options = append(options, deepseek.WithTimeout(c.Timeout))
	}
	if c.Path != "" {
		options = append(options, deepseek.WithPath(c.Path))
	}
	if c.ProxyUrl != "" {
		Url, err := url.Parse(c.ProxyUrl)
		if err != nil {
			panic(gone.ToError(err))
		}
		transport := http.Transport{
			Proxy: http.ProxyURL(Url),
		}
		client := http.Client{
			Transport: &transport,
		}
		options = append(options, deepseek.WithHTTPClient(&client))
	} else if httpDoer != nil {
		options = append(options, deepseek.WithHTTPClient(httpDoer))
	}
	return options
}

var clientMap sync.Map

func Load(loader gone.Loader) error {
	provider := gone.WrapFunctionProvider(func(tagConf string, param struct {
		configure gone.Configure `gone:"configure"`
		httpDoer  g.HTTPDoer     `gone:"deepseek-proxies" option:"allowNil"`
	}) (*deepseek.Client, error) {
		if value, ok := clientMap.Load(tagConf); ok {
			return value.(*deepseek.Client), nil
		}

		var prefix = "deepseek"
		_, keys := gone.TagStringParse(tagConf)
		if len(keys) > 0 && keys[0] != "" {
			prefix = strings.TrimSpace(keys[0])
		}
		var config Config
		err := param.configure.Get(prefix, &config, "")
		if err != nil {
			return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("get %s config err", prefix))
		}

		client, err := deepseek.NewClientWithOptions(config.AuthToken, config.ToDeepseekOptions(param.httpDoer)...)
		if err != nil {
			return nil, gone.ToErrorWithMsg(err, fmt.Sprintf("NewClientWithOptions by %s config err", prefix))
		}
		clientMap.Store(tagConf, client)
		return client, nil
	})

	load := gone.OnceLoad(func(loader gone.Loader) error {
		return loader.Load(provider)
	})
	return load(loader)
}
