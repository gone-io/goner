package kafka

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

type Conf struct {
	*sarama.Config
	Addrs   []string `json:"addrs"`
	GroupID string   `json:"groupID"`
}

var defaultConf = sarama.NewConfig()

func (c *Conf) ReadFromConfigure(name string, configure gone.Configure) ([]string, *sarama.Config) {
	if name == "" {
		name = "default"
	}
	confKey := fmt.Sprintf("kafka.%s", name)
	c.Config = sarama.NewConfig()
	_ = configure.Get(confKey, c, "")

	_ = mergo.Merge(c.Config, defaultConf)

	return c.Addrs, c.Config
}

func ProvideSyncProducer(tagConf string, param struct {
	configure  gone.Configure  `gone:"configure"`
	beforeStop gone.BeforeStop `gone:"*"`
	logger     gone.Logger     `gone:"*"`
}) (sarama.SyncProducer, error) {
	var conf Conf
	name, _ := gone.ParseGoneTag(tagConf)
	producer, err := sarama.NewSyncProducer(conf.ReadFromConfigure(name, param.configure))
	g.PanicIfErr(gone.ToErrorWithMsg(err, "create kafka sync producer failed"))
	param.beforeStop(func() {
		g.ErrorPrinter(param.logger, producer.Close(), "close kafka sync producer failed")
	})

	return producer, nil
}

func ProvideAsyncProducer(tagConf string, param struct {
	configure  gone.Configure  `gone:"configure"`
	beforeStop gone.BeforeStop `gone:"*"`
	logger     gone.Logger     `gone:"*"`
}) (sarama.AsyncProducer, error) {
	var conf Conf
	name, _ := gone.ParseGoneTag(tagConf)
	producer, err := sarama.NewAsyncProducer(conf.ReadFromConfigure(name, param.configure))
	g.PanicIfErr(gone.ToErrorWithMsg(err, "create kafka async producer failed"))

	param.beforeStop(func() {
		g.ErrorPrinter(param.logger, producer.Close(), "close kafka async producer failed")
	})

	return producer, nil
}

func ProvideConsumer(tagConf string, param struct {
	configure  gone.Configure  `gone:"configure"`
	beforeStop gone.BeforeStop `gone:"*"`
	logger     gone.Logger     `gone:"*"`
}) (sarama.Consumer, error) {
	var conf Conf
	name, _ := gone.ParseGoneTag(tagConf)
	consumer, err := sarama.NewConsumer(conf.ReadFromConfigure(name, param.configure))
	g.PanicIfErr(gone.ToErrorWithMsg(err, "create kafka consumer failed"))
	param.beforeStop(func() {
		g.ErrorPrinter(param.logger, consumer.Close(), "close kafka kafka consumer failed")
	})
	return consumer, nil
}

func ProvideConsumerGroup(tagConf string, param struct {
	configure  gone.Configure  `gone:"configure"`
	beforeStop gone.BeforeStop `gone:"*"`
	logger     gone.Logger     `gone:"*"`
}) (sarama.ConsumerGroup, error) {
	var conf Conf
	name, _ := gone.ParseGoneTag(tagConf)
	brokers, config := conf.ReadFromConfigure(name, param.configure)
	consumerGroup, err := sarama.NewConsumerGroup(brokers, conf.GroupID, config)
	g.PanicIfErr(gone.ToErrorWithMsg(err, "create kafka consumer group failed"))
	param.beforeStop(func() {
		g.ErrorPrinter(param.logger, consumerGroup.Close(), "close kafka consumer group failed")
	})
	return consumerGroup, nil
}
