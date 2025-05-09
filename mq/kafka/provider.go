package kafka

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gone-io/gone/v2"
)

type Conf struct {
	*sarama.Config
	Addrs []string `json:"addrs"`
}

var defaultConf = sarama.NewConfig()

func (c *Conf) ReadFromConfigure(name string, configure gone.Configure) ([]string, *sarama.Config) {
	if name == "" {
		name = "default"
	}
	confKey := fmt.Sprintf("kafka.%s", name)
	_ = configure.Get(confKey, c, "")

	_ = mergo.Merge(c.Config, defaultConf)

	return c.Addrs, c.Config
}

func ProvideSyncProducer(tagConf string, param struct {
	configure gone.Configure `gone:"configure"`
}) (sarama.SyncProducer, error) {
	var conf Conf
	name, _ := gone.ParseGoneTag(tagConf)
	producer, err := sarama.NewSyncProducer(conf.ReadFromConfigure(name, param.configure))
	return producer, gone.ToErrorWithMsg(err, "create kafka sync producer failed")
}

func ProvideAsyncProducer(tagConf string, param struct {
	configure gone.Configure `gone:"configure"`
}) (sarama.AsyncProducer, error) {
	var conf Conf
	name, _ := gone.ParseGoneTag(tagConf)
	producer, err := sarama.NewAsyncProducer(conf.ReadFromConfigure(name, param.configure))
	return producer, gone.ToErrorWithMsg(err, "create kafka async producer failed")
}

func ProvideConsumer(tagConf string, param struct {
	configure gone.Configure `gone:"configure"`
}) (sarama.Consumer, error) {
	var conf Conf
	name, _ := gone.ParseGoneTag(tagConf)
	consumer, err := sarama.NewConsumer(conf.ReadFromConfigure(name, param.configure))
	return consumer, gone.ToErrorWithMsg(err, "create kafka consumer failed")
}
