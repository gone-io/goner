package mqtt

import (
	"dario.cat/mergo"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gone-io/gone/v2"
)

var defaultOpts = mqtt.NewClientOptions()

var newClient = mqtt.NewClient

type Conf struct {
	mqtt.ClientOptions
	Brokers []string `json:"brokers"`
}

func (s *Conf) ToClientOptions() *mqtt.ClientOptions {
	for _, b := range s.Brokers {
		s.AddBroker(b)
	}
	_ = mergo.Merge(&s.ClientOptions, defaultOpts)
	return &s.ClientOptions
}

func ProvideClient(tagConf string, i struct {
	beforeStop gone.BeforeStop `gone:"*"`
	configure  gone.Configure  `gone:"configure"`
}) (mqtt.Client, error) {
	name, _ := gone.ParseGoneTag(tagConf)
	var conf Conf
	if name == "" {
		name = "default"
	}
	confKey := fmt.Sprintf("mqtt.%s", name)
	_ = i.configure.Get(confKey, &conf, "")

	client := newClient(conf.ToClientOptions())
	token := client.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return nil, gone.ToErrorWithMsg(err, "mqtt token failed")
	}
	i.beforeStop(func() {
		client.Disconnect(200)
	})
	return client, nil
}
