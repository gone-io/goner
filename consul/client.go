package consul

import (
	"github.com/gone-io/gone/v2"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
)

var client *api.Client

func GetDefaultConf(config *api.Config) *api.Config {
	if config.Address == "" {
		config.Address = "127.0.0.1:8500"
	}
	if config.Scheme == "" {
		config.Scheme = "http"
	}
	config.Transport = cleanhttp.DefaultPooledTransport()
	return config
}

func ProvideConsulClient(_ string, param struct {
	config *api.Config `gone:"config,consul"`
}) (*api.Client, error) {
	if client != nil {
		return client, nil
	}
	var err error
	client, err = api.NewClient(GetDefaultConf(param.config))
	if err != nil {
		return nil, gone.ToErrorWithMsg(err, "can not create consul client")
	}
	return client, nil
}
